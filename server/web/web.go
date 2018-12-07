package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/iyabchen/go-react-kv/server/data"
	"github.com/rs/cors"
)

// WebServer serves REST API to access data.
type WebServer struct {
	api    *API
	server *http.Server
	router *mux.Router
	addr   string
}

// Options for web server.
type Options struct {
	Addr    string
	Storage data.PairRepo
}

// NewWeb creates a new web server.
func NewWeb(opt *Options) (*WebServer, error) {
	api, err := NewAPI(opt.Storage)
	if err != nil {
		return nil, fmt.Errorf("Failed to init api: %s", err)
	}

	w := &WebServer{
		api:    api,
		addr:   opt.Addr,
		router: mux.NewRouter(),
	}
	w.registerAPI()
	return w, nil
}

// Run starts web server.
func (s *WebServer) Run() error {
	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		MaxAge: 300,
		AllowedHeaders: []string{
			"*",
		},
	})

	srv := &http.Server{
		Addr:         s.addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      corsOpts.Handler(handlers.LoggingHandler(os.Stdout, s.router)),
	}
	s.server = srv
	log.Printf("Listening on address %s", s.addr)
	return srv.ListenAndServe()
}

// Shutdown closes web server.
func (s *WebServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)

}

func (s *WebServer) registerAPI() {
	wrapper := func(f apiFunc) http.HandlerFunc {

		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			data, err := f(req)
			if err != nil {
				s.respondError(w, req, err)
			} else {
				if req.URL.String() == "/download" {
					w.Header().Set("Content-Description", "File Transfer")
					w.Header().Set("Content-Type", "application/zip")
					w.Header().Set("Content-Disposition", "attachment; filename=\"paris.json\"")
					w.Header().Set("Content-Transfer-Encoding", "binary")
				}

				s.respondOK(w, req, data)
			}

		})
	}

	s.router.HandleFunc("/pair", wrapper(s.api.getAll)).Methods("GET")
	s.router.HandleFunc("/pair/{id}", wrapper(s.api.getOne)).Methods("GET")
	s.router.HandleFunc("/pair", wrapper(s.api.create)).Methods("POST")
	s.router.HandleFunc("/pair/{id}", wrapper(s.api.deleteOne)).Methods("DELETE")
	s.router.HandleFunc("/pair/{id}", wrapper(s.api.update)).Methods("PUT")
	s.router.HandleFunc("/reset", wrapper(s.api.deleteAll)).Methods("GET")

	// the dir is relative path of the executable
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./client")))
}

type apiFunc func(c *http.Request) (data interface{}, err *apiError)

type apiError struct {
	httpCode int
	err      error
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s: %s", http.StatusText(e.httpCode), e.err)
}

func (s *WebServer) respondOK(w http.ResponseWriter, r *http.Request, data interface{}) {

	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			log.Printf("error marshaling response: %s", err)
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
		if _, err = w.Write(b); err != nil {
			log.Printf("error writing response: %s", err)
		}
	}
}

func (s *WebServer) respondError(w http.ResponseWriter, r *http.Request, apiErr *apiError) {
	type errResponse struct {
		Err string `json:"error"`
	}
	b, err := json.Marshal(&errResponse{
		Err: apiErr.Error(),
	})
	if err != nil {
		log.Printf("error marshalling json response: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.httpCode)
	if _, err := w.Write(b); err != nil {
		log.Printf("error writing response: %s", err)
	}
}
