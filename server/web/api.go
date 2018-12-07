package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/iyabchen/go-react-kv/server/data"
	"github.com/iyabchen/go-react-kv/server/model"
)

// API connects data repo.
type API struct {
	repo data.PairRepo
}

// NewAPI init API instance.
func NewAPI(repo data.PairRepo) (*API, error) {
	return &API{repo: repo}, nil
}

func (a *API) getAll(r *http.Request) (interface{}, *apiError) {
	data, err := a.repo.GetAll(context.Background())
	if err != nil {
		return nil, &apiError{http.StatusInternalServerError, err}
	}
	return data, nil

}

func (a *API) getOne(r *http.Request) (interface{}, *apiError) {
	params := mux.Vars(r)
	id := params["id"]
	data, err := a.repo.GetOne(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return nil, &apiError{
				httpCode: http.StatusNotFound,
				err:      fmt.Errorf("pair does not exist for id %s", id),
			}
		}
		return nil, &apiError{
			httpCode: http.StatusInternalServerError,
			err:      err,
		}
	}
	return data, nil

}

func (a *API) deleteAll(r *http.Request) (interface{}, *apiError) {
	err := a.repo.DeleteAll(context.Background())
	if err != nil {
		return nil, &apiError{http.StatusInternalServerError, err}
	}
	return nil, nil
}

func (a *API) deleteOne(r *http.Request) (interface{}, *apiError) {
	params := mux.Vars(r)
	id := params["id"]
	err := a.repo.DeleteOne(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return nil, &apiError{
				httpCode: http.StatusNotFound,
				err:      fmt.Errorf("pair does not exist for id %s", id),
			}
		}
		return nil, &apiError{
			httpCode: http.StatusInternalServerError,
			err:      err,
		}
	}
	return nil, nil
}

func validateInput(r *http.Request) (*model.Pair, error) {
	type tmp struct {
		Key   *string `json:"key"`
		Value *string `json:"value"`
	}
	var t tmp

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %s", err)
	}

	if t.Key == nil || t.Value == nil {
		return nil, fmt.Errorf("No content provided")

	}
	return model.NewPair(*t.Key, *t.Value)
}

func (a *API) create(r *http.Request) (interface{}, *apiError) {
	p, err := validateInput(r)
	if err != nil {
		return nil, &apiError{http.StatusBadRequest, err}
	}
	err = a.repo.Create(context.Background(), p)
	if err != nil {
		return nil, &apiError{http.StatusInternalServerError, err}
	}

	return p, nil
}

func (a *API) update(r *http.Request) (interface{}, *apiError) {
	p, err := validateInput(r)
	if err != nil {
		return nil, &apiError{http.StatusBadRequest, err}
	}
	params := mux.Vars(r)
	id := params["id"]
	err = a.repo.Update(context.Background(), id, p.Key, p.Value)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return nil, &apiError{
				httpCode: http.StatusNotFound,
				err:      fmt.Errorf("pair does not exist for id %s", id),
			}
		}
		return nil, &apiError{http.StatusInternalServerError, err}
	}

	return p, nil
}
