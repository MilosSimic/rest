package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"strconv"
)

func (ts *postServer) createPostHandler(w http.ResponseWriter, req *http.Request) {
	span := StartSpanFromRequest("cretePostHandler", ts.tracer, req)
	defer span.Finish()

	log.Printf("handling post create at %s\n", req.URL.Path)

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	ctx := ContextWithSpan(context.Background(), span)
	rt, err := decodeBody(ctx, req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreatePost(ctx, rt.Title, rt.Text, rt.Tags)
	renderJSON(ctx, w, ResponseId{Id: id})
}

func (ts *postServer) getAllPostsHandler(w http.ResponseWriter, req *http.Request) {
	span := StartSpanFromRequest("getAllPostsHandler", ts.tracer, req)
	defer span.Finish()

	log.Printf("handling get all posts at %s\n", req.URL.Path)

	ctx := ContextWithSpan(context.Background(), span)
	allTasks := ts.store.GetAllPosts(ctx)
	renderJSON(ctx, w, allTasks)
}

func (ts *postServer) getPostHandler(w http.ResponseWriter, req *http.Request) {
	span := StartSpanFromRequest("getPostHandler", ts.tracer, req)
	defer span.Finish()

	log.Printf("handling get post at %s\n", req.URL.Path)

	ctx := ContextWithSpan(context.Background(), span)
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.GetPost(ctx, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(ctx, w, task)
}

func (ts *postServer) deletePostHandler(w http.ResponseWriter, req *http.Request) {
	span := StartSpanFromRequest("deletePostHandler", ts.tracer, req)
	defer span.Finish()

	log.Printf("handling delete post at %s\n", req.URL.Path)

	ctx := ContextWithSpan(context.Background(), span)
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := ts.store.DeletePost(ctx, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
