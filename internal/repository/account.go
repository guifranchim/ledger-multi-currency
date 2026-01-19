package repository

import (
	"errors"
	"ledger-multi-currency/internal/domain"
	"sync"
)

type AccountRepository interface {
	Create(account *domain.Account) error
	GetByID(id string) (*domain.Account, error)
	List() ([]*domain.Account, error)
	Update(account *domain.Account) error
	Delete(id string) error
}

type accountRepositoryImpl struct {
	mu       sync.RWMutex
	accounts map[string]*domain.Account
}

func NewAccountRepository() AccountRepository {
	return &accountRepositoryImpl{
		accounts: make(map[string]*domain.Account),
	}
}

func (r *accountRepositoryImpl) Create(account *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[account.ID]; exists {
		return errors.New("account: já existe conta com esse ID")
	}

	r.accounts[account.ID] = account
	return nil
}

func (r *accountRepositoryImpl) GetByID(id string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[id]
	if !exists {
		return nil, errors.New("account: não encontrada")
	}
	return account, nil
}

func (r *accountRepositoryImpl) List() ([]*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	accounts := make([]*domain.Account, 0, len(r.accounts))
	for _, acc := range r.accounts {
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (r *accountRepositoryImpl) Update(account *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[account.ID]; !exists {
		return errors.New("account: não encontrada")
	}

	r.accounts[account.ID] = account
	return nil
}

func (r *accountRepositoryImpl) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[id]; !exists {
		return errors.New("account: não encontrada")
	}

	delete(r.accounts, id)
	return nil
}
