package repository

import (
	"errors"
	"ledger-multi-currency/internal/domain"
	"sync"
)

type JournalRepository interface {
	Create(entry *domain.JournalEntry) error
	GetByID(id string) (*domain.JournalEntry, error)
	List() ([]*domain.JournalEntry, error)
	Post(id string) error
}

type journalRepositoryImpl struct {
	mu      sync.RWMutex
	entries map[string]*domain.JournalEntry
}

func NewJournalRepository() JournalRepository {
	return &journalRepositoryImpl{
		entries: make(map[string]*domain.JournalEntry),
	}
}

func (r *journalRepositoryImpl) Create(entry *domain.JournalEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entries[entry.ID]; exists {
		return errors.New("journal: já existe lançamento com esse ID")
	}

	r.entries[entry.ID] = entry
	return nil
}

func (r *journalRepositoryImpl) GetByID(id string) (*domain.JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.entries[id]
	if !exists {
		return nil, errors.New("journal: lançamento não encontrado")
	}
	return entry, nil
}

func (r *journalRepositoryImpl) List() ([]*domain.JournalEntry, error) {
	r.mu.RLock()
	defer r.mu.Unlock()

	entries := make([]*domain.JournalEntry, 0, len(r.entries))
	for _, entry := range r.entries {
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *journalRepositoryImpl) Post(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, exists := r.entries[id]
	if !exists {
		return errors.New("journal: lançamento não encontrado")
	}

	if entry.Status != "DRAFT" {
		return errors.New("journal: apenas lançamentos em DRAFT podem ser postados")
	}

	if !entry.IsBalanced() {
		return errors.New("journal: lançamento não balanceado")
	}

	entry.Status = "POSTED"
	return nil
}
