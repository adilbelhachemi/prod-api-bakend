package storage

import (
	"pratbacknd/internal/types"
)

type UpdateProductInput struct {
	ProductId        string      `json:"productId"`
	Name             string      `json:"name"`
	Image            string      `json:"image"`
	ShortDescription string      `json:"shortDescription"`
	Description      string      `json:"description"`
	PriceVATExcluded types.Money `json:"priceVATExcluded"`
	VAT              types.Money `json:"vat"`
	TotalPrice       types.Money `json:"totalPrice"`
}

type Storage interface {
	Products() ([]types.Product, error)
	GetProductById(productID string) (types.Product, error)
	CreateProduct(p types.Product) error
	UpdateProduct(input UpdateProductInput) error

	Categories() ([]types.Category, error)
	CreateCategory(c types.Category) error

	UpdateInventory(productId string, delta int) error

	CreateCart(cart types.Cart, userId string) error
	GetCart(userID string) (types.Cart, error)

	CreateOrUpdateCart(userID string, productID string, delta int) (types.Cart, error)
}
