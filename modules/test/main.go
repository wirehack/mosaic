package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
)

func Wire(di *dig.Container) {
	di.Invoke(func(router chi.Router) {
		router.Get("/id", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello world!"))
		})
	})
}
