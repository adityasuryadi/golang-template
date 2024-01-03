package repository

import (
	"order-service/internal/entity"
	"order-service/internal/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(tx *gorm.DB, order *entity.Order) error
	FindAllByProductId(tx *gorm.DB)
	Search(db *gorm.DB, request *model.SearchOrderRequest) ([]entity.Order, int64, error)
	FilterOrder(request *model.SearchOrderRequest) func(tx *gorm.DB) *gorm.DB
}
