package handler

import (
	"encoding/json"
	"ledger-multi-currency/internal/domain"
	"ledger-multi-currency/internal/service"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	service *service.AccountService
}

func NewAccountHandler(service *service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

type CreateAccountRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Currency string `json:"currency"`
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "erro ao decodificar request", http.StatusBadRequest)
		return
	}

	account, err := h.service.CreateAccount(req.ID, req.Name, req.Type, req.Currency)
	if err != nil {
		slog.Error("erro ao criar conta", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	account, err := h.service.GetAccount(id)
	if err != nil {
		slog.Error("erro ao buscar conta", "id", id, "err", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListAccounts()
	if err != nil {
		slog.Error("erro ao listar contas", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if accounts == nil {
		accounts = make([]*domain.Account, 0)
	}
	json.NewEncoder(w).Encode(accounts)
}

type GetBalanceResponse struct {
	Balances map[string]float64 `json:"balances"`
}

func (h *AccountHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	balances, err := h.service.GetAccountBalance(id)
	if err != nil {
		slog.Error("erro ao buscar saldo", "id", id, "err", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := GetBalanceResponse{
		Balances: make(map[string]float64),
	}
	for currency, minor := range balances {
		response.Balances[currency] = float64(minor) / 100.0
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AccountHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeactivateAccount(id); err != nil {
		slog.Error("erro ao desativar conta", "id", id, "err", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
