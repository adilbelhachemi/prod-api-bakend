package server

import (
	"os"
	"pratbacknd/internal/category"
	"pratbacknd/internal/product"

	"github.com/Rhymond/go-money"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	App  *fiber.App
	port string
}

type Config struct {
	Port string
}

func New(config Config) (*Server, error) {
	fiberApp := fiber.New()

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173/, https://master.d14f8mlnk4lkw2.amplifyapp.com/",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	s := &Server{App: fiberApp, port: config.Port}

	fiberApp.Get("/products", s.Products())
	fiberApp.Get("/categories", s.Categories())

	return s, nil
}

func (s *Server) Run() error {
	// Listen and Server in 0.0.0.0:$PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = s.port
	}

	return s.App.Listen(":" + port)
}

func (s *Server) Categories() fiber.Handler {
	return func(c *fiber.Ctx) error {
		categories := []category.Category{
			category.Category{
				ID:          "11",
				Name:        "Test",
				Description: "this the first category",
			},
			category.Category{
				ID:          "12",
				Name:        "Test 2",
				Description: "This is the 2nd categoty",
			},
		}
		return c.Status(fiber.StatusOK).JSON(categories)
	}
}

func (s *Server) Products() fiber.Handler {
	return func(c *fiber.Ctx) error {
		categories := []product.Product{
			product.Product{
				ID:               "42",
				Name:             "Test",
				Description:      "This is my product",
				PriceVATExcluded: money.New(100, "EUR"),
				VAT:              money.New(200, "EUR")},
			product.Product{
				ID:               "33",
				Name:             "Test 2",
				Description:      "This is my 2nd product",
				PriceVATExcluded: money.New(70, "EUR"),
				VAT:              money.New(140, "EUR")},
		}
		return c.Status(fiber.StatusOK).JSON(categories)
	}
}
