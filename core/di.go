package core

import (
	"c"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DI interface {
	DB() *pgxpool.Pool
	Router() chi.Router
}

type DII struct {
	DI
	router chi.Router
}

func (dim *DII) DB() *pgxpool.Pool  { return c.DB() }
func (dim *DII) Router() chi.Router { return dim.router }

func NewDI(router chi.Router) DI {
	return &DII{
		router: router,
	}
}
