package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server, err := NewPostServer()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer server.GetCloser()

	router.HandleFunc("/post/", server.createPostHandler).Methods("POST")
	router.HandleFunc("/post/", server.getAllPostsHandler).Methods("GET")
	router.HandleFunc("/post/{id:[0-9]+}/", server.getPostHandler).Methods("GET")
	router.HandleFunc("/post/{id:[0-9]+}/", server.deletePostHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
