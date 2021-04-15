package main

import (
	ps "example.com/internal/poststore"
	opentracing "github.com/opentracing/opentracing-go"
	"io"
)

type postServer struct {
	store  *ps.PostStore
	tracer opentracing.Tracer
	closer io.Closer
}

const name = "post_service"

func NewPostServer() (*postServer, error) {
	store, err := ps.New()
	if err != nil {
		return nil, err
	}

	tracer, closer := Init(name)
	return &postServer{
		store:  store,
		tracer: tracer,
		closer: closer,
	}, nil
}

func (s *postServer) GetTracer() opentracing.Tracer {
	return s.tracer
}

func (s *postServer) GetCloser() io.Closer {
	return s.closer
}
