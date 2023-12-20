package repository

import (
	"order-service/internal/entity"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(tx *gorm.DB, order *entity.Order) error
	FindAllByProductId(tx *gorm.DB)
}
