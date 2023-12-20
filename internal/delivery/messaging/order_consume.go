package messaging

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OrderConsumer struct {
	Log *logrus.Logger
}

func NewOrderConsumer(log *logrus.Logger) *OrderConsumer {
	return &OrderConsumer{
		Log: log,
	}
}

type Product struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
	Name        string    `gorm:"column:name"`
	Price       float64   `gorm:"column:price"`
	Qty         int       `gorm:"column:qty"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (c OrderConsumer) Consume(message []byte) error {
	product := new(Product)
	err := json.Unmarshal(message, product)
	if err != nil {
		log.Fatal("failed unmarshal")
		return err
	}
	c.Log.Infof("processDeliveries deliveryTag% v", product)
	return nil
}
