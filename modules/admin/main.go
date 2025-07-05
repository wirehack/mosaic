package main

import (
	"mosaic/core"
	"net/http"
)

type proxy struct{}

func (proxy) Meta() *core.ModuleInfo {
	return &core.ModuleInfo{
		Slug:        "admin",
		Description: "Admin module",
		Version:     "0.0.0",
	}
}

func Wire(di core.DI) any {

	di.Router().Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ADMIN"))
	})

	return &proxy{}
}
