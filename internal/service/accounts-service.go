package service

import (
	"ledger-multi-currency/internal/domain"
	"ledger-multi-currency/internal/repository"
)

type AccountsService struct {
	repo repository.AccountsRepository
}

func NewAccountsService(repo repository.AccountsRepository) *AccountsService {
	return &AccountsService{
		repo: repo,
	}
}

func (s *AccountsService) GetAccountByID(accountID string) (*domain.Accounts, error) {

	return s.repo.GetByID(accountID)
}

func (s *AccountsService) CreateAccount(account *domain.Accounts) error {

	return s.repo.Create(account)
}

func (s *AccountsService) ListAccounts() ([]*domain.Accounts, error) {
	return s.repo.List()
}
