package types

import (
	"fmt"

	"github.com/Rhymond/go-money"
)

type Cart struct {
	ID           string
	CurrencyCode string
	Items        map[string]Item
	Version      uint
}

type Item struct {
	ID               string
	ShortDescription string
	Quantity         uint8
	UnitPriceVATExc  *money.Money
	VAT              *money.Money
	UnitPriceVATInc  *money.Money
}

func (c Cart) TotalPriceVATInc() (*money.Money, error) {
	totalPrice := money.New(0, c.CurrencyCode)
	for _, item := range c.Items {
		itemPrice := item.UnitPriceVATInc.Multiply(int64(item.Quantity))
		var err error
		totalPrice, err = totalPrice.Add(itemPrice)
		if err != nil {
			return nil, fmt.Errorf("error - add item price to total price: %w", err)
		}
	}
	return totalPrice, nil
}

func (c *Cart) UpsertItem(productID string, delta int) error {
	item, found := c.Items[productID]
	if !found {
		// item is not in the cart, we have to add it
		if delta <= 0 {
			return fmt.Errorf("error - item not found, delta is less or equal than zero: (delta = %d)", delta)
		}
		c.Items[productID] = Item{
			ID:       productID,
			Quantity: uint8(delta),
		}
	} else {
		// a product with this id is already in the cart
		newQuantity := int(item.Quantity) + delta
		if newQuantity < 0 {
			return fmt.Errorf("error - new quantity cannot be less than zero")
		} else if newQuantity > 0 {
			item.Quantity = uint8(newQuantity)
			c.Items[productID] = item
		} else {
			// equal to zero
			// we need to remove from the cart
			delete(c.Items, productID)
		}
	}

	return nil
}
