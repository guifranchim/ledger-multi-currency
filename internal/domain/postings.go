package domain

type Posting struct {
	ID             string
	EntryID        string
	AccountID      string
	AmountMinor    int64
	Currency       string
	AmountBRLMinor int64
	FXRateID       *string
}
