package main

import (
	"mosaic/core"

	"go.uber.org/dig"
)

type proxy struct{}

func (proxy) Meta() *core.ModuleInfo {
	return &core.ModuleInfo{
		Slug:        "mail",
		Description: "Mail send module",
		Version:     "0.0.0",
	}
}

func Wire(di *dig.Scope) any {
	return &proxy{}
}
