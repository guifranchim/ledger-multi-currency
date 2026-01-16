package routes

import (
	"fmt"
	"ledger-multi-currency/internal/handler"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(r chi.Router, accountsHandler *handler.AccountsHandler) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.HealthHandler)

	setupAPIV1Routes(r, accountsHandler)

	printRoutes(r)
}

func setupAPIV1Routes(r chi.Router, accountsHandler *handler.AccountsHandler) {
	r.Route("/api/v1", func(apiV1 chi.Router) {

		apiV1.Route("/accounts", func(accounts chi.Router) {
			accounts.Get("/", accountsHandler.List)
			accounts.Post("/", accountsHandler.Create)
			accounts.Get("/{id}", accountsHandler.GetByID)
		})
	})
}

func printRoutes(r chi.Router) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		slog.Info(fmt.Sprintf("[%s] %s", method, route))
		return nil
	}

	if mux, ok := r.(*chi.Mux); ok {
		chi.Walk(mux, walkFunc)
	}
}
