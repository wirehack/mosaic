package main

import (
	"c"
	"mosaic/core"
	"mosaic/route"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	di.Router().Post("/", p.create)
	di.Router().Get("/{id}", p.me)
	di.Router().Post("/auth", p.noop)
	di.Router().Put("/acm/{resource}/{permission}", p.noop)
	di.Router().Get("/acm/{resource}/{permission}", p.noop)

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

func main() {

	root := chi.NewMux()

	router := route.CreateMainRouter()
	root.Mount("/v1", route.MountRoutes(router))

	di := core.NewDI(router)

	Wire(di)

	c.PrintRoutes(root)

	if err := http.ListenAndServe(":8080", root); err != nil {
		panic(err)
	}

}
