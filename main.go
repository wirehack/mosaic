package main

import (
	"c"
	"mosaic/core"
	"mosaic/route"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	root := chi.NewMux()
	router := route.CreateMainRouter()
	root.Mount("/v1", route.MountRoutes(router))

	router.Route("/modules", func(r chi.Router) {
		scope := c.DI().Scope("modules")
		scope.Provide(func() chi.Router { return r })
		if _, err := core.RegisterModules(scope); err != nil {
			panic(err)
		}
	})

	if err := http.ListenAndServe(":8080", root); err != nil {
		panic(err)
	}
}
