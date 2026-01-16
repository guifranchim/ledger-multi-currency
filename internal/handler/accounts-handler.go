package handler

import (
	"encoding/json"
	"ledger-multi-currency/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AccountsHandler struct {
	service *service.AccountsService
}

func NewAccountsHandler(service *service.AccountsService) *AccountsHandler {
	return &AccountsHandler{
		service: service,
	}
}

func (h *AccountsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")

	account, err := h.service.GetAccountByID(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (h *AccountsHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *AccountsHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListAccounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}
