package repository

import (
	"order-service/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	Repository[entity.Product]
	Log *logrus.Logger
}

// FindProductsById implements Productrepository.
func (r *ProductRepositoryImpl) FindProductsById(tx *gorm.DB, Id []string) ([]entity.Product, error) {
	var product []entity.Product
	err := tx.Where("id IN ?", Id).Find(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func NewProductRepository(log *logrus.Logger) Productrepository {
	return &ProductRepositoryImpl{
		Log: log,
	}
}
