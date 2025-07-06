package main

import (
	"context"
	"mosaic/core"
	"mosaic/types"
)

type proxy struct {
	di core.DI
}

func (p *proxy) Send(recipient types.SendRecipient, subject, template string, params *types.D) error {
	p.di.DB().Exec(context.Background(), "SELECT 1")
	return nil
}

func (p *proxy) Meta() *core.ModuleInfo {
	return &core.ModuleInfo{
		Slug:        "mail",
		Description: "Mail send module",
		Version:     "0.0.0",
	}
}

func Wire(di core.DI) any {
	return &proxy{di}
}
