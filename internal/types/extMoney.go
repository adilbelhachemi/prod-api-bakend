package types

type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Display  string `json:"display"`
}
