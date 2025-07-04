package main

import (
	"mosaic/core"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
)

type proxy struct{}

func (proxy) Meta() *core.ModuleInfo {
	return &core.ModuleInfo{
		Slug:        "id",
		Description: "id",
		Version:     "0.0.1",
	}
}

func noop(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func Wire(di *dig.Scope) any {

	di.Invoke(func(router chi.Router) {
		router.Post("/id", noop)
		router.Get("/id/{id}", noop)
		router.Put("/id/acm/{resource}/{permission}", noop)
		router.Get("/id/acm/{resource}/{permission}", noop)
	})

	return &proxy{}
}
