package main

import (
	"context"
	"order-service/internal/config"
	"order-service/internal/delivery/messaging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	channel, _ := config.NewRabbitMqChannell(viperConfig, log)
	orderHandler := messaging.NewOrderConsumer(log)

	consumeOrderCfg := messaging.ConsumerConfig{
		Exchange:       "product.created",
		QueueName:      "product.create",
		RoutingKey:     "create",
		ConsumerTag:    "",
		BindingKey:     "create",
		WorkerPoolSize: 5,
	}

	go messaging.Consume(consumeOrderCfg, channel, orderHandler.Consume)

	log.Info("Worker is running")

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	stop := false
	for !stop {
		select {
		case s := <-terminateSignals:
			log.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
			cancel()
			stop = true
		}
	}

	time.Sleep(5 * time.Second)

}
