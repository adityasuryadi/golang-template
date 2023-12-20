package config

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRabbitMqChannell(viper *viper.Viper, log *logrus.Logger) (*amqp.Channel, error) {
	connAddr := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		viper.GetString("rabbitmq.username"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetInt("rabbitmq.port"))
	conn, err := amqp.Dial(connAddr)
	if err != nil {
		log.WithError(err).Error("failed to connect rabbitmq")
	}
	ch, err := conn.Channel()
	if err != nil {
		log.WithError(err).Error("failed to create channel")
	}
	return ch, nil
}
