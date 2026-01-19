package repository

import (
	"errors"
	"ledger-multi-currency/internal/domain"
	"sync"
)

type FXRateRepository interface {
	Create(rate *domain.FXRate) error
	GetLatest(from, to string) (*domain.FXRate, error)
	List() ([]*domain.FXRate, error)
}

type fxrateRepositoryImpl struct {
	mu    sync.RWMutex
	rates map[string]*domain.FXRate
	index map[string]string
}

func NewFXRateRepository() FXRateRepository {
	return &fxrateRepositoryImpl{
		rates: make(map[string]*domain.FXRate),
		index: make(map[string]string),
	}
}

func (r *fxrateRepositoryImpl) Create(rate *domain.FXRate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rates[rate.ID]; exists {
		return errors.New("fxrate: já existe taxa com esse ID")
	}

	r.rates[rate.ID] = rate
	key := rate.FromCurrency + "_" + rate.ToCurrency
	r.index[key] = rate.ID

	return nil
}

func (r *fxrateRepositoryImpl) GetLatest(from, to string) (*domain.FXRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := from + "_" + to
	id, exists := r.index[key]
	if !exists {
		return nil, errors.New("fxrate: taxa não encontrada")
	}

	rate, _ := r.rates[id]
	return rate, nil
}

func (r *fxrateRepositoryImpl) List() ([]*domain.FXRate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rates := make([]*domain.FXRate, 0, len(r.rates))
	for _, rate := range r.rates {
		rates = append(rates, rate)
	}
	return rates, nil
}
