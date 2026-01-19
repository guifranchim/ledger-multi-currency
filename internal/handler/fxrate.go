package handler

import (
	"encoding/json"
	"ledger-multi-currency/internal/domain"
	"ledger-multi-currency/internal/service"
	"log/slog"
	"net/http"
)

type FXRateHandler struct {
	service *service.FXRateService
}

func NewFXRateHandler(service *service.FXRateService) *FXRateHandler {
	return &FXRateHandler{service: service}
}

type RegisterRateRequest struct {
	ID           string `json:"id"`
	FromCurrency string `json:"fromCurrency"`
	ToCurrency   string `json:"toCurrency"`
	RateScaled   int64  `json:"rateScaled"`
	Scale        int    `json:"scale"`
}

func (h *FXRateHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "erro ao decodificar request", http.StatusBadRequest)
		return
	}

	rate, err := h.service.RegisterRate(
		req.ID,
		req.FromCurrency,
		req.ToCurrency,
		req.RateScaled,
		req.Scale,
	)
	if err != nil {
		slog.Error("erro ao registrar taxa", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rate)
}

func (h *FXRateHandler) List(w http.ResponseWriter, r *http.Request) {
	rates, err := h.service.ListRates()
	if err != nil {
		slog.Error("erro ao listar taxas", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if rates == nil {
		rates = make([]*domain.FXRate, 0)
	}
	json.NewEncoder(w).Encode(rates)
}

type ConvertRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

type ConvertResponse struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Original  int64  `json:"original"`
	Converted int64  `json:"converted"`
	Rate      struct {
		Numerator   int64 `json:"numerator"`
		Denominator int64 `json:"denominator"`
	} `json:"rate"`
}

func (h *FXRateHandler) Convert(w http.ResponseWriter, r *http.Request) {
	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "erro ao decodificar request", http.StatusBadRequest)
		return
	}

	converted, err := h.service.Convert(req.From, req.To, req.Amount)
	if err != nil {
		slog.Error("erro ao converter", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rate, _ := h.service.GetLatestRate(req.From, req.To)

	response := ConvertResponse{
		From:      req.From,
		To:        req.To,
		Original:  req.Amount,
		Converted: converted,
	}
	if rate != nil {
		response.Rate.Numerator = rate.RateScaled
		response.Rate.Denominator = int64(rate.Scale)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
