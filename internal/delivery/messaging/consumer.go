package messaging

import (
	"context"
	"errors"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
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

type Consumer struct {
	channell *amqp.Channel
	pCfg     *ConsumerConfig
	Log      *logrus.Logger
}

type ConsumerConfig struct {
	Exchange       string
	QueueName      string
	RoutingKey     string
	ConsumerTag    string
	BindingKey     string
	WorkerPoolSize int
}
type ConsumeHandler func(message []byte) error

func (c *Consumer) SetupExchangeAndQueueConsumer(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch := c.channell

	err := ch.ExchangeDeclare(
		exchangeName,
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments)
	)

	if err != nil {
		return nil, errors.New("Error ch.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)

	if err != nil {
		return nil, errors.New("Error ch.QueueDeclare")
	}

	fmt.Printf("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		exchangeName,
		bindingKey)

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.New("Error ch.QueueBind")
	}

	fmt.Printf("Queue bound to exchange, starting to consume from queue, consumerTag: %v", consumerTag)

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		return nil, errors.New("Error  ch.Qos")
	}

	return ch, nil
}

func (c *Consumer) worker(ctx context.Context, messages <-chan amqp.Delivery) {
	for delivery := range messages {
		fmt.Printf("processDeliveries deliveryTag% v", delivery.DeliveryTag)
		// unmarshal here
	}
	fmt.Println("Deliveries channel closed")
}

func (c *Consumer) StartConsumer(consumerCfg ConsumerConfig, handler ConsumeHandler) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := c.SetupExchangeAndQueueConsumer(consumerCfg.Exchange, consumerCfg.QueueName, consumerCfg.BindingKey, consumerCfg.ConsumerTag)
	if err != nil {
		return errors.New(err.Error())
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		consumerCfg.QueueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return errors.New("Consume")
	}

	for i := 0; i < consumerCfg.WorkerPoolSize; i++ {
		go c.worker(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	fmt.Printf("ch.NotifyClose: %v", chanErr)
	return chanErr
}

// test consume

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
