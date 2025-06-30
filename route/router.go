package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CreateMainRouter() chi.Router {
	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(middleware.RequestID)
	return mux
}

func MountRoutes(router chi.Router) chi.Router {
	router.Get("/id", CreateID)
	return router
}
