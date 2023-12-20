package messaging

import (
	"order-service/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type OrderPublisher struct {
	Publisher[*model.OrderEvent]
}

func NewOrderProducer(ch *amqp.Channel, pubConfig *PublisherConfig, log *logrus.Logger) *OrderPublisher {
	return &OrderPublisher{
		Publisher: Publisher[*model.OrderEvent]{
			channell: ch,
			pCfg:     pubConfig,
			Log:      log,
		},
	}
}
