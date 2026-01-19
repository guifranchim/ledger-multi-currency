package domain

import (
	"testing"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("ACC001", "Caixa", "ASSET", "BRL")

	if err != nil {
		t.Errorf("NewAccount() error = %v", err)
	}

	if acc.ID != "ACC001" {
		t.Errorf("NewAccount() ID = %s, want ACC001", acc.ID)
	}
}

func TestJournalBalance(t *testing.T) {
	entry, _ := NewJournalEntry("J001", "Venda", "NF001")

	p1, _ := NewPosting("P001", "ACC001", "USD", 10000, true)
	entry.AddPosting(p1)

	p2, _ := NewPosting("P002", "ACC002", "USD", 10000, false)
	entry.AddPosting(p2)

	if !entry.IsBalanced() {
		t.Error("IsBalanced() = false, want true")
	}
}

func TestFXRateConvert(t *testing.T) {
	rate, _ := NewFXRate("FX001", "USD", "BRL", 180000, 1000000)
	result := rate.Convert(10000)

	if result != 1800 {
		t.Errorf("Convert() = %d, want 1800", result)
	}
}
