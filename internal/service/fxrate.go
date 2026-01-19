package service

import (
	"fmt"
	"ledger-multi-currency/internal/domain"
	"ledger-multi-currency/internal/repository"
)

type FXRateService struct {
	fxrateRepo repository.FXRateRepository
}

func NewFXRateService(fxrateRepo repository.FXRateRepository) *FXRateService {
	return &FXRateService{
		fxrateRepo: fxrateRepo,
	}
}

func (s *FXRateService) RegisterRate(id, fromCurrency, toCurrency string, rateScaled int64, scale int) (*domain.FXRate, error) {
	rate, err := domain.NewFXRate(id, fromCurrency, toCurrency, rateScaled, scale)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar taxa: %w", err)
	}

	if err := s.fxrateRepo.Create(rate); err != nil {
		return nil, fmt.Errorf("erro ao persistir taxa: %w", err)
	}

	return rate, nil
}

func (s *FXRateService) GetLatestRate(from, to string) (*domain.FXRate, error) {
	rate, err := s.fxrateRepo.GetLatest(from, to)
	if err != nil {
		return nil, fmt.Errorf("taxa não encontrada de %s para %s: %w", from, to, err)
	}
	return rate, nil
}

func (s *FXRateService) Convert(from, to string, amountMinor int64) (int64, error) {
	rate, err := s.GetLatestRate(from, to)
	if err != nil {
		return 0, err
	}

	converted := rate.Convert(amountMinor)
	return converted, nil
}

func (s *FXRateService) ListRates() ([]*domain.FXRate, error) {
	rates, err := s.fxrateRepo.List()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar taxas: %w", err)
	}
	return rates, nil
}
