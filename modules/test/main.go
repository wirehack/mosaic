package main

import (
	"core"
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
)

//go:embed ui/dist/*
var UI embed.FS

func Wire(di *dig.Container) {
	di.Invoke(func(router chi.Router) {
		router.Get("/id", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello world!"))
		})
	})
}

func Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Slug:        "test",
		Description: "test",
		Version:     "1.1.1",
	}
}
