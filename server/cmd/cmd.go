package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/iyabchen/go-react-kv/server/data"
	"github.com/iyabchen/go-react-kv/server/web"
)

func main() {
	addr := ":8080"
	if os.Getenv("PORT") != "" {
		addr = fmt.Sprintf(":" + os.Getenv("PORT"))
	}

	ds, err := data.NewMem()
	if err != nil {
		log.Fatalf("Failed to create data storage: %s", err)
	}

	srv, err := web.NewWeb(&web.Options{Addr: addr, Storage: ds})
	if err != nil {
		log.Fatalf("Failed to create web server: %s", err)
	}

	go func() {
		// service connections
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
