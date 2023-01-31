package server

import (
	"net/http"
	"pratbacknd/internal/category"
	"pratbacknd/internal/product"

	"github.com/Rhymond/go-money"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router         *chi.Mux
	port           string
	allowedOrigins string
}

type Config struct {
	Port           string
	AllowedOrigins string
}

func New(config Config) (*Server, error) {
	r := chi.NewRouter()
	s := &Server{Router: r, port: config.Port, allowedOrigins: config.AllowedOrigins}

	r.Use(s.enableCORS)

	r.Get("/products", s.Products)
	r.Get("/categories", s.Categories)

	return s, nil
}

func (s *Server) Categories(w http.ResponseWriter, r *http.Request) {
	categories := []category.Category{
		{
			ID:          "11",
			Name:        "Test",
			Description: "this the first category",
		},
		{
			ID:          "12",
			Name:        "Test 2",
			Description: "This is the 2nd categoty",
		},
	}
	s.writeJSON(w, http.StatusOK, categories)
}

func (s *Server) Products(w http.ResponseWriter, r *http.Request) {
	products := []product.Product{
		{
			ID:               "42",
			Name:             "Test",
			Description:      "This is my product",
			PriceVATExcluded: money.New(100, "EUR"),
			VAT:              money.New(200, "EUR")},
		{
			ID:               "33",
			Name:             "Test 2",
			Description:      "This is my 2nd product",
			PriceVATExcluded: money.New(70, "EUR"),
			VAT:              money.New(140, "EUR")},
	}
	s.writeJSON(w, http.StatusOK, products)
}
