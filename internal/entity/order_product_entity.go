package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderProduct struct {
	gorm.Model
	Id           uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
	OrderId      uuid.UUID `gorm:"column:order_id;type:uuid;"`
	ProductId    uuid.UUID `gorm:"column:product_id;type:uuid"`
	ProductName  string    `gorm:"column:product_name"`
	ProductPrice float64   `gorm:"column:product_price"`
	Qty          int64     `gorm:"column:qty"`
	TotalPrice   float64   `gorm:"total_price"`
	CreatedAt    int64     `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt    int64     `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	Order        Order
	Product      Product `gorm:"foreignKey:ProductId;"`
}

func (entity *OrderProduct) TableName() string {
	return "order_products"
}

func (entity *OrderProduct) BeforeCreate(db *gorm.DB) error {
	entity.Id = uuid.New()
	return nil
}
