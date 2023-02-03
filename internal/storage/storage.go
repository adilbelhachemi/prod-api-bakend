package storage

import (
	"pratbacknd/internal/category"
	"pratbacknd/internal/product"
)

type Storage interface {
	Products() ([]product.Product, error)
	Categories() ([]category.Category, error)

	CreateProduct(p product.Product) error
	CreateCategory(c category.Category) error
}
