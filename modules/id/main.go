package main

import (
	"mosaic/core"
	"net/http"
)

type proxy struct {
	di core.DI
}

func (p proxy) Meta() *core.ModuleInfo {
	return &core.ModuleInfo{
		Slug:        "id",
		Description: "id",
		Version:     "0.0.1",
	}
}

func Wire(di core.DI) any {

	p := proxy{di}

	di.Router().Post("/id", p.create)
	di.Router().Get("/id/{id}", p.me)
	di.Router().Put("/id/acm/{resource}/{permission}", p.noop)
	di.Router().Get("/id/acm/{resource}/{permission}", p.noop)

	return p
}

func (p proxy) create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func (p proxy) me(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func (p proxy) noop(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}
