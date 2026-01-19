package main

import (
	"ledger-multi-currency/internal/handler"
	"ledger-multi-currency/internal/repository"
	"ledger-multi-currency/internal/routes"
	"ledger-multi-currency/internal/service"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	accountRepo := repository.NewAccountRepository()
	journalRepo := repository.NewJournalRepository()
	fxrateRepo := repository.NewFXRateRepository()

	accountService := service.NewAccountService(accountRepo, journalRepo)
	journalService := service.NewJournalService(journalRepo, accountRepo)
	fxrateService := service.NewFXRateService(fxrateRepo)

	accountHandler := handler.NewAccountHandler(accountService)
	journalHandler := handler.NewJournalHandler(journalService)
	fxrateHandler := handler.NewFXRateHandler(fxrateService)

	r := chi.NewRouter()
	routes.Setup(r, accountHandler, journalHandler, fxrateHandler)

	slog.Info("Starting Ledger API on :3000...")
	if err := http.ListenAndServe(":3000", r); err != nil {
		slog.Error("Server error", "err", err)
	}
}
