package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)
	server, err := NewPostServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer server.CloseTracer()
	defer server.CloseDB()

	router.HandleFunc("/post/", count(server.createPostHandler)).Methods("POST")
	router.HandleFunc("/post/", count(server.getAllPostsHandler)).Methods("GET")
	router.HandleFunc("/post/{id:[0-9]+}/", count(server.getPostHandler)).Methods("GET")
	router.HandleFunc("/post/{id:[0-9]+}/", count(server.deletePostHandler)).Methods("DELETE")
	router.Path("/metrics").Handler(metricsHandler())

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}
