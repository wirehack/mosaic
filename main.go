package main

import (
	"mosaic/core"
	"mosaic/route"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	root := chi.NewMux()

	router := route.CreateMainRouter()
	root.Mount("/v1", route.MountRoutes(router))

	if _, err := core.RegisterModules(core.DI()); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":8080", root); err != nil {
		panic(err)
	}
}
