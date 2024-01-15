package database

import "github.com/kleytonsolinho/golang-client-server-api/internal/entity"

type CurrencyInterface interface {
	Create(currency *entity.CurrencyPrice) error
}
