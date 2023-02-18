package utils

import (
	"pratbacknd/internal/types"
	reflect "reflect"
)

type peoductUpdateInput struct {
	productId        string
	name             string
	image            string
	shortDescription string
	description      string
	priceVatExcluded types.Money
	vat              types.Money
	totalPrice       types.Money
}

func ValidateProductUpdateInput(personMap map[string]interface{}) bool {
	productInput := peoductUpdateInput{}
	productValue := reflect.ValueOf(&productInput).Elem()

	for key := range personMap {
		field := productValue.FieldByName(key)
		if !field.IsValid() {
			return false
		}
	}
	return true
}
