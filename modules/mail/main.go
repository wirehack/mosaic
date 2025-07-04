package main

import (
	"context"
	"mosaic/core"
	"mosaic/types"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/dig"
)

type proxy struct {
	di *dig.Scope
	types.MailModuleProxy
}

func (p *proxy) Send(recipient types.SendRecipient, subject, template string, params *types.D) error {
	p.di.Invoke(func(db func() *pgxpool.Pool) error {
		db().Exec(context.Background(), "SELECT 1")
		return nil
	})
	return nil
}

func (p *proxy) Meta() *core.ModuleInfo {
	return &core.ModuleInfo{
		Slug:        "mail",
		Description: "Mail send module",
		Version:     "0.0.0",
	}
}

func Wire(di *dig.Scope) any {
	return &proxy{
		di: di,
	}
}
