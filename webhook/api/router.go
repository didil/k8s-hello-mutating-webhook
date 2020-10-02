package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// BuildRouter builds the router
func BuildRouter(app *App) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/mutate", app.HandleMutate)

	return r
}
