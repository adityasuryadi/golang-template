package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	Id           uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
	ProductId    uuid.UUID `gorm:"column:product_id;type:uuid"`
	UserId       uuid.UUID `gorm:"column:user_id;type:uuid"`
	UserName     string    `gorm:"column:user_name"`
	ProductName  string    `gorm:"column:product_name"`
	ProductPrice float64   `gorm:"column:product_price"`
	Qty          int64     `gorm:"column:qty"`
	CreatedAt    int64     `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt    int64     `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (entity *Order) TableName() string {
	return "orders"
}

func (entity *Order) BeforeCreate(db *gorm.DB) error {
	entity.Id = uuid.New()
	return nil
}
