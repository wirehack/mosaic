package main

import (
	"mosaic/core"
	"net/http"
)

type proxy struct{}

func (proxy) Meta() *core.ModuleInfo {
	return &core.ModuleInfo{
		Slug:        "id",
		Description: "id",
		Version:     "0.0.1",
	}
}

func Wire(di core.DI) any {

	di.Router().Post("/id", create)
	di.Router().Get("/id/{id}", me)
	di.Router().Put("/id/acm/{resource}/{permission}", noop)
	di.Router().Get("/id/acm/{resource}/{permission}", noop)

	return &proxy{}
}

func create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func me(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func noop(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}
