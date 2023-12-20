package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Publisher[T any] struct {
	channell *amqp.Channel
	pCfg     *PublisherConfig
	Log      *logrus.Logger
}

type PublisherConfig struct {
	Exchange    string
	QueueName   string
	RoutingKey  string
	ConsumerTag string
}

func (p *Publisher[T]) SetupExchangeAndQueuePublisher() {
	fmt.Println("declare exchange")
	ch := p.channell
	err := ch.ExchangeDeclare(
		p.pCfg.Exchange, // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments)
	)
	p.Log.WithError(err).Error("Failed to declare an exchange")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	q, err := ch.QueueDeclare(
		p.pCfg.QueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	p.Log.WithError(err).Error("Failed to declare a queue")

	err = ch.QueueBind(q.Name, p.pCfg.RoutingKey, p.pCfg.Exchange, false, nil)
	p.Log.WithError(err).Error("Failed to declare a queue")
}

func (p *Publisher[T]) CloseChannel() {
	if err := p.channell.Close(); err != nil {
		p.Log.WithError(err).Error("Failed to close channel")
	}
}

func (p *Publisher[T]) Publish(event T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	value, err := json.Marshal(event)
	if err != nil {
		p.Log.WithError(err).Error("failed to marshal event")
	}
	defer cancel()
	p.channell.PublishWithContext(ctx, p.pCfg.Exchange, p.pCfg.RoutingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        value,
	})
	log.Printf(" [x] Sent %s", value)
}
