package main

import (
	"database/sql"
	"ledger-multi-currency/internal/handler"
	"ledger-multi-currency/internal/repository"
	"ledger-multi-currency/internal/routes"
	"ledger-multi-currency/internal/service"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	var db *sql.DB

	accountsRepo := repository.NewAccountsRepository(db)
	accountsService := service.NewAccountsService(accountsRepo)

	accountsHandler := handler.NewAccountsHandler(accountsService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	routes.SetupRoutes(r, accountsHandler)

	slog.Info("Server started on :3000")
	http.ListenAndServe(":3000", r)
}
