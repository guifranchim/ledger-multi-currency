package service

import (
	"testing"

	"ledger-multi-currency/internal/repository"
)

func TestAccountServiceCreate(t *testing.T) {
	repo := repository.NewAccountRepository()
	journalRepo := repository.NewJournalRepository()
	service := NewAccountService(repo, journalRepo)

	account, err := service.CreateAccount("ACC001", "Caixa", "ASSET", "BRL")
	if err != nil {
		t.Fatalf("CreateAccount() error = %v", err)
	}

	if account.ID != "ACC001" {
		t.Errorf("CreateAccount() ID = %s, want ACC001", account.ID)
	}
}

func TestFXRateServiceConvert(t *testing.T) {
	fxrateRepo := repository.NewFXRateRepository()
	fxrateService := NewFXRateService(fxrateRepo)

	fxrateService.RegisterRate("FX001", "USD", "BRL", 180000, 1000000)
	result, err := fxrateService.Convert("USD", "BRL", 10000)

	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	if result != 1800 {
		t.Errorf("Convert() = %d, want 1800", result)
	}
}
