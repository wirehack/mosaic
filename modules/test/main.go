package main

import (
	"embed"
	"mosaic/core"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
)

//go:embed ui/dist/**
var UIFS embed.FS

type proxy struct{}

func (proxy) UI() *embed.FS {
	return &UIFS
}

func (proxy) Meta() *core.ModuleInfo {
	return &core.ModuleInfo{
		Slug:        "test",
		Description: "test",
		Version:     "1.1.1",
	}
}

func (proxy) Sum(a, b int) int { return a + b }

func Wire(di *dig.Scope) any {

	di.Invoke(func(router chi.Router) {
		router.Get("/id", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello world!"))
		})
	})

	return &proxy{}
}
