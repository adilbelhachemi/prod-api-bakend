package storage

import (
	"pratbacknd/internal/types"
)

type UpdateProductInput struct {
	ProductId        string
	Name             string
	Image            string
	ShortDescription string
	Description      string
	PriceVATExcluded types.Amount
	VAT              types.Amount
	TotalPrice       types.Amount
}

type Storage interface {
	Products() ([]types.Product, error)
	CreateProduct(p types.Product) error
	UpdateProduct(input UpdateProductInput) error

	Categories() ([]types.Category, error)
	CreateCategory(c types.Category) error

	UpdateInventory(productId string, delta int) error

	CreateCart(cart types.Cart, userId string) error
	GetCart(userID string) (types.Cart, error)
}
