package repository

import (
	"order-service/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	Repository[entity.Order]
	Log *logrus.Logger
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

func NewOrderRepository(log *logrus.Logger) OrderRepository {
	return &OrderRepositoryImpl{
		Log: log,
	}
}
