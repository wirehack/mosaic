package main

import (
	"mosaic/core"
	"net/http"
)

type proxy struct {
	di core.DI
}

func (proxy) Meta() *core.ModuleInfo {
	return &core.ModuleInfo{
		Slug:        "admin",
		Description: "Admin module",
		Version:     "0.0.0",
	}
}

func (p proxy) AdminFiles(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ADMIN"))
}

func Wire(di core.DI) any {

	p := &proxy{di}

	di.Router().Get("/admin/ui", p.AdminFiles)

	return p
}
