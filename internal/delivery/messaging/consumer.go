package messaging

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	publishMandatory = false
	publishImmediate = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

type ConsumerConfig struct {
	Exchange       string
	QueueName      string
	RoutingKey     string
	ConsumerTag    string
	BindingKey     string
	WorkerPoolSize int
}

type ConsumeHandler func(message []byte) error

func Consume(cfg ConsumerConfig, ch *amqp.Channel, handler ConsumeHandler) {

	err := ch.ExchangeDeclare(
		cfg.Exchange,
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments)
	)
	if err != nil {
		log.Fatalf("Error ch.ExchangeDeclare : %v", err)
	}

	queue, err := ch.QueueDeclare(
		cfg.QueueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)

	if err != nil {
		log.Fatalf("Error ch.QueueDeclare : %v", err)
	}

	fmt.Printf("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		cfg.Exchange,
		cfg.BindingKey)

	err = ch.QueueBind(
		queue.Name,
		cfg.BindingKey,
		cfg.Exchange,
		queueNoWait,
		nil,
	)
	if err != nil {
		log.Fatalf("Error ch.QueueBind : %v", err)
	}

	fmt.Printf("Queue bound to exchange, starting to consume from queue, consumerTag: %v", cfg.ConsumerTag)

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		log.Fatalf("Error  ch.Qos : %v", err)
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer ch.Close()

	deliveries, err := ch.Consume(
		cfg.QueueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Error  Consume : %v", err)
	}
	for i := 0; i < cfg.WorkerPoolSize; i++ {
		go func() {
			for delivery := range deliveries {
				err := handler(delivery.Body)
				if err != nil {
					delivery.Nack(true, true)
				}
				delivery.Ack(true)
			}
			fmt.Println("Deliveries channel closed")
		}()
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	fmt.Printf("ch.NotifyClose: %v", chanErr)
}
