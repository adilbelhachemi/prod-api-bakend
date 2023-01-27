package product

import "github.com/Rhymond/go-money"

type Product struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	Description      string       `json:"description"`
	PriceVATExcluded *money.Money `json:"priceVATExcluded"`
	VAT              *money.Money `json:"vat"`
}
