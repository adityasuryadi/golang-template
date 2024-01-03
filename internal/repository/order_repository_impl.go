package repository

import (
	"order-service/internal/entity"
	"order-service/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	Repository[entity.Order]
	Log *logrus.Logger
}

// Search implements OrderRepository.
func (r *OrderRepositoryImpl) Search(db *gorm.DB, request *model.SearchOrderRequest) ([]entity.Order, int64, error) {
	var orders []entity.Order
	err := db.Scopes(r.FilterOrder(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Preload("OrderProducts").Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Order{}).Scopes(r.FilterOrder(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// FindAllByProductId implements OrderRepository.
func (r *OrderRepositoryImpl) FindAllByProductId(tx *gorm.DB) {
	panic("unimplemented")
}

func (r *OrderRepositoryImpl) Create(tx *gorm.DB, order *entity.Order) error {
	err := r.Repository.Create(tx, order)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepositoryImpl) FilterOrder(request *model.SearchOrderRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx
	}
}

func NewOrderRepository(log *logrus.Logger) OrderRepository {
	return &OrderRepositoryImpl{
		Log: log,
	}
}
