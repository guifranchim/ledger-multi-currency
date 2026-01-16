package repository

import (
	"database/sql"
	"ledger-multi-currency/internal/domain"
)

type AccountsRepository interface {
	GetByID(accountID string) (*domain.Accounts, error)
	Create(account *domain.Accounts) error
	List() ([]*domain.Accounts, error)
}

type accountsRepositoryImpl struct {
	db *sql.DB
}

func NewAccountsRepository(db *sql.DB) AccountsRepository {
	return &accountsRepositoryImpl{
		db: db,
	}
}

func (r *accountsRepositoryImpl) GetByID(accountID string) (*domain.Accounts, error) {

	return &domain.Accounts{
		ID:       accountID,
		Name:     "Account from DB",
		Type:     "Asset",
		Currency: "USD",
	}, nil
}

func (r *accountsRepositoryImpl) Create(account *domain.Accounts) error {

	return nil
}

func (r *accountsRepositoryImpl) List() ([]*domain.Accounts, error) {

	return []*domain.Accounts{}, nil
}
