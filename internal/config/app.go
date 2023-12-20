package config

import (
	controller "order-service/internal/delivery/http"
	"order-service/internal/delivery/http/route"
	"order-service/internal/repository"
	"order-service/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// userRepository := repository.NewUserRepository(config.Log)
	// contactRepository := repository.NewContactRepository(config.Log)
	// addressRepository := repository.NewAddressRepository(config.Log)

	orderRepository := repository.NewOrderRepository(config.Log)
	orderUsecase := usecase.NewOrderUsecase(config.DB, config.Log, orderRepository)
	orderController := controller.NewOrderController(orderUsecase, config.Log)
	routeConfig := route.RouteConfig{
		App:             config.App,
		OrderController: orderController,
	}
	routeConfig.Setup()
}
