package domain

import (
	"errors"
	"fmt"
	"time"
)

type Account struct {
	ID        string
	Name      string
	Type      string
	Currency  string
	CreatedAt time.Time
	Active    bool
}

func NewAccount(id, name, accountType, currency string) (*Account, error) {
	if id == "" || name == "" || accountType == "" || currency == "" {
		return nil, errors.New("account: campos obrigatórios")
	}

	validTypes := map[string]bool{
		"ASSET": true, "LIABILITY": true, "EQUITY": true,
		"INCOME": true, "EXPENSE": true,
	}
	if !validTypes[accountType] {
		return nil, fmt.Errorf("account: tipo inválido %s", accountType)
	}

	return &Account{
		ID:        id,
		Name:      name,
		Type:      accountType,
		Currency:  currency,
		CreatedAt: time.Now(),
		Active:    true,
	}, nil
}

type JournalEntry struct {
	ID          string
	Description string
	Reference   string
	Postings    []*Posting
	CreatedAt   time.Time
	Status      string
}

func NewJournalEntry(id, description, reference string) (*JournalEntry, error) {
	if id == "" || description == "" {
		return nil, errors.New("journal: id e description são obrigatórios")
	}

	return &JournalEntry{
		ID:        id,
		Description: description,
		Reference:   reference,
		Postings:    make([]*Posting, 0),
		CreatedAt:   time.Now(),
		Status:      "DRAFT",
	}, nil
}

func (je *JournalEntry) AddPosting(p *Posting) error {
	if je.Status == "POSTED" {
		return errors.New("journal: não pode modificar lançamento já postado")
	}
	je.Postings = append(je.Postings, p)
	return nil
}

func (je *JournalEntry) IsBalanced() bool {
	// Agrupa saldos por moeda
	balances := make(map[string]int64)
	for _, p := range je.Postings {
		balances[p.Currency] += p.GetAmountSigned()
	}

	for _, balance := range balances {
		if balance != 0 {
			return false
		}
	}
	return true
}

type Posting struct {
	ID          string
	AccountID   string
	Currency    string
	AmountMinor int64
	Debit       bool
	CreatedAt   time.Time
}

func NewPosting(id, accountID, currency string, amountMinor int64, debit bool) (*Posting, error) {
	if id == "" || accountID == "" || currency == "" {
		return nil, errors.New("posting: campos obrigatórios")
	}
	if amountMinor <= 0 {
		return nil, errors.New("posting: valor deve ser positivo")
	}

	return &Posting{
		ID:          id,
		AccountID:   accountID,
		Currency:    currency,
		AmountMinor: amountMinor,
		Debit:       debit,
		CreatedAt:   time.Now(),
	}, nil
}

func (p *Posting) GetAmountSigned() int64 {
	if p.Debit {
		return p.AmountMinor
	}
	return -p.AmountMinor
}

func (p *Posting) GetAmountDecimal() float64 {
	return float64(p.AmountMinor) / 100.0
}

type FXRate struct {
	ID           string
	FromCurrency string
	ToCurrency   string
	RateScaled   int64
	Scale        int
	CreatedAt    time.Time
}

func NewFXRate(id, from, to string, rateScaled int64, scale int) (*FXRate, error) {
	if id == "" || from == "" || to == "" {
		return nil, errors.New("fxrate: campos obrigatórios")
	}
	if from == to {
		return nil, errors.New("fxrate: moedas não podem ser iguais")
	}
	if rateScaled <= 0 || scale <= 0 {
		return nil, errors.New("fxrate: taxa e escala devem ser positivas")
	}

	return &FXRate{
		ID:           id,
		FromCurrency: from,
		ToCurrency:   to,
		RateScaled:   rateScaled,
		Scale:        scale,
		CreatedAt:    time.Now(),
	}, nil
}

func (f *FXRate) GetRate() float64 {
	return float64(f.RateScaled) / float64(f.Scale)
}

func (f *FXRate) Convert(amountMinor int64) int64 {
	return (amountMinor * f.RateScaled) / int64(f.Scale)
}
