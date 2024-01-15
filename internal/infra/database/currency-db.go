package database

import (
	"context"

	"github.com/kleytonsolinho/golang-client-server-api/internal/entity"
	"gorm.io/gorm"
)

type CurrencyRepository struct {
	DB *gorm.DB
}

func NewCurrencyDb(db *gorm.DB) *CurrencyRepository {
	return &CurrencyRepository{
		DB: db,
	}
}

func (c *CurrencyRepository) Create(ctx context.Context, currencyPrice *entity.CurrencyPrice) error {
	return c.DB.WithContext(ctx).Create(currencyPrice).Error
}
