package main

import (
	"core"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
)

func main() {

	mux := chi.NewMux()

	di := dig.New()
	di.Provide(func() chi.Router { return mux })

	if err := core.RegisterModules(di); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
