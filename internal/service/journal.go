package service

import (
	"errors"
	"fmt"
	"ledger-multi-currency/internal/domain"
	"ledger-multi-currency/internal/repository"
)

type JournalService struct {
	journalRepo repository.JournalRepository
	accountRepo repository.AccountRepository
}

func NewJournalService(journalRepo repository.JournalRepository, accountRepo repository.AccountRepository) *JournalService {
	return &JournalService{
		journalRepo: journalRepo,
		accountRepo: accountRepo,
	}
}

func (s *JournalService) CreateJournalEntry(id, description, reference string) (*domain.JournalEntry, error) {
	entry, err := domain.NewJournalEntry(id, description, reference)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar lançamento: %w", err)
	}

	if err := s.journalRepo.Create(entry); err != nil {
		return nil, fmt.Errorf("erro ao persistir lançamento: %w", err)
	}

	return entry, nil
}

func (s *JournalService) AddPosting(entryID, postingID, accountID, currency string, amountMinor int64, debit bool) error {
	entry, err := s.journalRepo.GetByID(entryID)
	if err != nil {
		return fmt.Errorf("lançamento não encontrado: %w", err)
	}

	if _, err := s.accountRepo.GetByID(accountID); err != nil {
		return fmt.Errorf("conta não encontrada: %w", err)
	}

	posting, err := domain.NewPosting(postingID, accountID, currency, amountMinor, debit)
	if err != nil {
		return fmt.Errorf("erro ao criar posting: %w", err)
	}

	if err := entry.AddPosting(posting); err != nil {
		return fmt.Errorf("erro ao adicionar posting: %w", err)
	}

	return nil
}

func (s *JournalService) PostJournalEntry(entryID string) error {
	entry, err := s.journalRepo.GetByID(entryID)
	if err != nil {
		return fmt.Errorf("lançamento não encontrado: %w", err)
	}

	if len(entry.Postings) == 0 {
		return errors.New("lançamento sem postings")
	}

	if err := s.journalRepo.Post(entryID); err != nil {
		return fmt.Errorf("erro ao postar lançamento: %w", err)
	}

	return nil
}

func (s *JournalService) GetJournalEntry(id string) (*domain.JournalEntry, error) {
	entry, err := s.journalRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("lançamento não encontrado: %w", err)
	}
	return entry, nil
}

func (s *JournalService) ListJournalEntries() ([]*domain.JournalEntry, error) {
	entries, err := s.journalRepo.List()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar lançamentos: %w", err)
	}
	return entries, nil
}
