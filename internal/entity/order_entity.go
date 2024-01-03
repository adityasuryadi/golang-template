package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Id               uuid.UUID      `gorm:"primaryKey;type:uuid;" column:"id"`
	OrderNo          string         `gorm:"column:order_no"`
	OrderNoCounter   int64          `gorm:"column:order_no_counter"`
	UserId           uuid.UUID      `gorm:"column:user_id;type:uuid"`
	OrderName        string         `gorm:"column:order_name"`
	ShippmentAddress string         `gorm:"column:shippment_address"`
	CreatedAt        int64          `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt        int64          `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	OrderProducts    []OrderProduct `gorm:"foreignKey:OrderId;references:Id;"`
	TotalPrice       float64        `gorm:"column:total_price"`
}

func (entity *Order) TableName() string {
	return "orders"
}

func (entity *Order) BeforeCreate(db *gorm.DB) error {
	entity.Id = uuid.New()
	return nil
}
