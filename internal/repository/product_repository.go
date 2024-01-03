package repository

import (
	"order-service/internal/entity"

	"gorm.io/gorm"
)

type Productrepository interface {
	FindProductsById(tx *gorm.DB, Id []string) ([]entity.Product, error)
}
