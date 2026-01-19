package routes

import (
	"fmt"
	"ledger-multi-currency/internal/handler"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Setup(
	r *chi.Mux,
	accountHandler *handler.AccountHandler,
	journalHandler *handler.JournalHandler,
	fxrateHandler *handler.FXRateHandler,
) {

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ledger API - Multi Currency"))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(apiV1 chi.Router) {

		apiV1.Route("/accounts", func(accounts chi.Router) {
			accounts.Post("/", accountHandler.Create)
			accounts.Get("/", accountHandler.List)
			accounts.Get("/{id}", accountHandler.GetByID)
			accounts.Get("/{id}/balance", accountHandler.GetBalance)
			accounts.Delete("/{id}", accountHandler.Deactivate)
		})

		apiV1.Route("/journals", func(journals chi.Router) {
			journals.Post("/", journalHandler.Create)
			journals.Get("/", journalHandler.List)
			journals.Get("/{id}", journalHandler.GetByID)
			journals.Post("/{id}/postings", journalHandler.AddPosting)
			journals.Post("/{id}/post", journalHandler.Post)
		})

		apiV1.Route("/rates", func(rates chi.Router) {
			rates.Post("/", fxrateHandler.Register)
			rates.Get("/", fxrateHandler.List)
			rates.Post("/convert", fxrateHandler.Convert)
		})
	})

	printRoutes(r)
}

func printRoutes(r *chi.Mux) {
	slog.Info("=== Rotas configuradas ===")
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		slog.Info(fmt.Sprintf("%s %s", method, route))
		return nil
	}

	chi.Walk(r, walkFunc)
	slog.Info("===========================")
}
