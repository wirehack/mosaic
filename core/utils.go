package core

import (
	"c"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/dig"
)

var __di = dig.New()

func XDI() *dig.Scope {
	scope := __di.Scope("root")
	scope.Provide(func() func() *pgxpool.Pool {
		return c.DB
	})
	return scope
}

func Module[T any](name string) (proxy T, exists bool) {
	XDI().Invoke(func(mp *ModuleProxy) {
		if ret, has := mp.Get(name); has {
			exists = true
			proxy = ret.(T)
		}
	})
	return
}
