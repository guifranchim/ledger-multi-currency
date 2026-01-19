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
	ID           string  `json:"id"`
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Rate         float64 `json:"rate"`
}

func (h *FXRateHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "erro ao decodificar request", http.StatusBadRequest)
		return
	}

	rateScaled := int64(req.Rate * 1000000)
	scale := 1000000

	rate, err := h.service.RegisterRate(
		req.ID,
		req.FromCurrency,
		req.ToCurrency,
		rateScaled,
		scale,
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
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type ConvertResponse struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Original  float64 `json:"original"`
	Converted float64 `json:"converted"`
	Rate      float64 `json:"rate"`
}

func (h *FXRateHandler) Convert(w http.ResponseWriter, r *http.Request) {
	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "erro ao decodificar request", http.StatusBadRequest)
		return
	}

	amountMinor := int64(req.Amount * 100)

	converted, err := h.service.Convert(req.From, req.To, amountMinor)
	if err != nil {
		slog.Error("erro ao converter", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rate, _ := h.service.GetLatestRate(req.From, req.To)
	rateValue := 0.0
	if rate != nil {
		rateValue = rate.GetRate()
	}

	response := ConvertResponse{
		From:      req.From,
		To:        req.To,
		Original:  req.Amount,
		Converted: float64(converted) / 100.0,
		Rate:      rateValue,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
