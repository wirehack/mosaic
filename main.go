package main

import (
	"c"
	"core"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
)

func main() {

	di := dig.New()

	router := core.CreateMainRouter()

	router.Route("/modules", func(r chi.Router) {
		scope := di.Scope("modules")
		scope.Provide(func() chi.Router { return r })
		if err := core.RegisterModules(scope); err != nil {
			panic(err)
		}
	})

	root := chi.NewMux()
	root.Mount("/v1", router)

	c.PrintRoutes(root)

	if err := http.ListenAndServe(":8080", root); err != nil {
		panic(err)
	}
}
