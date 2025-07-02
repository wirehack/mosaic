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
		Slug:        "admin",
		Description: "Admin module",
		Version:     "0.0.0",
	}
}

func Wire(di *dig.Scope) any {

	di.Invoke(func(router chi.Router) {
		router.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ADMIN"))
		})
	})

	return &proxy{}
}
