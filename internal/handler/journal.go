package handler

import (
	"encoding/json"
	"ledger-multi-currency/internal/domain"
	"ledger-multi-currency/internal/service"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type JournalHandler struct {
	service *service.JournalService
}

func NewJournalHandler(service *service.JournalService) *JournalHandler {
	return &JournalHandler{service: service}
}

type CreateJournalRequest struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Reference   string `json:"reference"`
}

func (h *JournalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateJournalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "erro ao decodificar request", http.StatusBadRequest)
		return
	}

	entry, err := h.service.CreateJournalEntry(req.ID, req.Description, req.Reference)
	if err != nil {
		slog.Error("erro ao criar lançamento", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

func (h *JournalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	entry, err := h.service.GetJournalEntry(id)
	if err != nil {
		slog.Error("erro ao buscar lançamento", "id", id, "err", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

func (h *JournalHandler) List(w http.ResponseWriter, r *http.Request) {
	entries, err := h.service.ListJournalEntries()
	if err != nil {
		slog.Error("erro ao listar lançamentos", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if entries == nil {
		entries = make([]*domain.JournalEntry, 0)
	}
	json.NewEncoder(w).Encode(entries)
}

type AddPostingRequest struct {
	PostingID string  `json:"postingId"`
	AccountID string  `json:"accountId"`
	Currency  string  `json:"currency"`
	Amount    float64 `json:"amount"`
	Debit     bool    `json:"debit"`
}

func (h *JournalHandler) AddPosting(w http.ResponseWriter, r *http.Request) {
	entryID := chi.URLParam(r, "id")

	var req AddPostingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "erro ao decodificar request", http.StatusBadRequest)
		return
	}

	amountMinor := int64(req.Amount * 100)

	err := h.service.AddPosting(
		entryID,
		req.PostingID,
		req.AccountID,
		req.Currency,
		amountMinor,
		req.Debit,
	)
	if err != nil {
		slog.Error("erro ao adicionar posting", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *JournalHandler) Post(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.PostJournalEntry(id); err != nil {
		slog.Error("erro ao postar lançamento", "id", id, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
