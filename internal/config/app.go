package config

import (
	controller "order-service/internal/delivery/http"
	"order-service/internal/delivery/http/route"
	"order-service/internal/pkg"
	"order-service/internal/repository"
	"order-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB             *gorm.DB
	App            *fiber.App
	Log            *logrus.Logger
	Validate       *pkg.Validation
	Config         *viper.Viper
	JaegerExporter *otlptrace.Exporter
}

func Bootstrap(config *BootstrapConfig) {
	// userRepository := repository.NewUserRepository(config.Log)
	// contactRepository := repository.NewContactRepository(config.Log)
	// addressRepository := repository.NewAddressRepository(config.Log)
	JaegerTracer := pkg.NewJaegerTracer(config.JaegerExporter)

	orderRepository := repository.NewOrderRepository(config.Log)
	productRepository := repository.NewProductRepository(config.Log)
	orderUsecase := usecase.NewOrderUsecase(config.DB, config.Log, orderRepository, productRepository, JaegerTracer, config.Validate)
	orderController := controller.NewOrderController(orderUsecase, config.Log, config.Validate)
	routeConfig := route.RouteConfig{
		App:             config.App,
		OrderController: orderController,
	}
	routeConfig.Setup()
}
