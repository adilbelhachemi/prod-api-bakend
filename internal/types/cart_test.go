package types

import (
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/stretchr/testify/assert"
)

func TestCart_TotalPriceVATInc(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		//  given
		cart := Cart{
			ID:           "1",
			CurrencyCode: "EUR",
			Items: map[string]Item{
				"43": {
					ID:               "11",
					ShortDescription: "a pair of socks",
					Quantity:         1,
					UnitPriceVATExc:  money.New(50, "EUR"),
					VAT:              money.New(50, "EUR"),
					UnitPriceVATInc:  money.New(100, "EUR"),
				},
			},
		}

		// when
		actualTotalPrice, err := cart.TotalPriceVATInc()

		// then
		assert.NoError(t, err, "error computing total price VAT included")

		expectedTotalPrice := money.New(100, "EUR")
		assert.Equal(t, expectedTotalPrice, actualTotalPrice)
	})

	t.Run("nominal with quantity greater than 1 and 2 items", func(t *testing.T) {
		// given
		items := map[string]Item{
			"42": {
				ID:               "42",
				ShortDescription: "A pair of socks",
				UnitPriceVATInc:  money.New(100, "EUR"),
				UnitPriceVATExc:  money.New(50, "EUR"),
				VAT:              money.New(50, "EUR"),
				Quantity:         1,
			},
			"43": {
				ID:               "43",
				ShortDescription: "A T-Shirt with a small gopher",
				UnitPriceVATInc:  money.New(3480, "EUR"),
				UnitPriceVATExc:  money.New(2900, "EUR"),
				VAT:              money.New(580, "EUR"),
				Quantity:         2,
			},
		}
		cart := Cart{
			ID:           "42",
			CurrencyCode: "EUR",
			Items:        items,
		}

		// when
		actualTotalPrice, err := cart.TotalPriceVATInc()

		// then
		assert.NoError(t, err, "error computing total price VAT included")
		expectedPriceVATINC := money.New(7060, "EUR")
		assert.Equal(t, expectedPriceVATINC, actualTotalPrice)
	})

	t.Run("error case different currencies", func(t *testing.T) {
		items := map[string]Item{
			"42": {
				ID:               "42",
				ShortDescription: "A pair of socks",
				UnitPriceVATInc:  money.New(100, "USD"),
				UnitPriceVATExc:  money.New(50, "USD"),
				VAT:              money.New(50, "USD"),
				Quantity:         1,
			},
		}
		cart := Cart{
			ID:           "42",
			CurrencyCode: "EUR",
			Items:        items,
		}

		// WHEN
		_, err := cart.TotalPriceVATInc()

		// THEN
		assert.Error(t, err, "when I add an item with a currency X to a basket of currency Y the method TotalPriceVATInc should fail")
	})
}
