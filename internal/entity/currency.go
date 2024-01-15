package entity

import "github.com/google/uuid"

type CurrencyPrice struct {
	ID    string `json:"id"`
	Code  string `json:"code"`
	Price string `json:"price"`
}

func NewCurrencyPrice(price string) (*CurrencyPrice, error) {
	currency := &CurrencyPrice{
		ID:    uuid.New().String(),
		Code:  "USDBRL",
		Price: price,
	}

	return currency, nil
}
