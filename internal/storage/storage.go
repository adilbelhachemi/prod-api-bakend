package storage

import (
	"pratbacknd/internal/category"
	"pratbacknd/internal/product"
)

type UpdateProductInput struct {
	ProductId        string
	Name             string
	Image            string
	ShortDescription string
	Description      string
	PriceVATExcluded product.Amount
	VAT              product.Amount
	TotalPrice       product.Amount
}

type Storage interface {
	Products() ([]product.Product, error)
	CreateProduct(p product.Product) error
	UpdateProduct(input UpdateProductInput) error

	Categories() ([]category.Category, error)
	CreateCategory(c category.Category) error

	UpdateInventory(productId string, delta int) error
}
