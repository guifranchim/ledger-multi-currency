package domain

type FXRate struct {
	ID         string
	Base       string
	Quote      string
	RateScaled int64
	Scale      int
	Timestamp  string
}
