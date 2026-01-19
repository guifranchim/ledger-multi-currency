package service

import (
	"fmt"
	"ledger-multi-currency/internal/domain"
	"ledger-multi-currency/internal/repository"
)

type AccountService struct {
	accountRepo repository.AccountRepository
	journalRepo repository.JournalRepository
}

func NewAccountService(accountRepo repository.AccountRepository, journalRepo repository.JournalRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		journalRepo: journalRepo,
	}
}

func (s *AccountService) CreateAccount(id, name, accountType, currency string) (*domain.Account, error) {
	account, err := domain.NewAccount(id, name, accountType, currency)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar conta: %w", err)
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, fmt.Errorf("erro ao persistir conta: %w", err)
	}

	return account, nil
}

func (s *AccountService) GetAccount(id string) (*domain.Account, error) {
	account, err := s.accountRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta: %w", err)
	}
	return account, nil
}

func (s *AccountService) ListAccounts() ([]*domain.Account, error) {
	accounts, err := s.accountRepo.List()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas: %w", err)
	}
	return accounts, nil
}

func (s *AccountService) GetAccountBalance(accountID string) (map[string]int64, error) {
	if _, err := s.accountRepo.GetByID(accountID); err != nil {
		return nil, fmt.Errorf("conta não encontrada: %w", err)
	}

	entries, err := s.journalRepo.List()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar lançamentos: %w", err)
	}

	balances := make(map[string]int64)
	for _, entry := range entries {
		if entry.Status != "POSTED" {
			continue
		}

		for _, posting := range entry.Postings {
			if posting.AccountID == accountID {
				balances[posting.Currency] += posting.GetAmountSigned()
			}
		}
	}

	return balances, nil
}

func (s *AccountService) DeactivateAccount(id string) error {
	account, err := s.GetAccount(id)
	if err != nil {
		return err
	}

	account.Active = false
	return s.accountRepo.Update(account)
}
