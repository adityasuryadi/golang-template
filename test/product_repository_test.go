package test

import (
	"fmt"
	"order-service/internal/config"
	"order-service/internal/repository"
	"testing"

	"github.com/spf13/viper"
)

func TestFindProductsByMultiplleId(t *testing.T) {
	viperConfig := viper.New()

	viperConfig.SetConfigName("config-test")
	viperConfig.SetConfigType("json")
	viperConfig.AddConfigPath("./../")
	viperConfig.AddConfigPath("./")
	err := viperConfig.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	productRepository := repository.NewProductRepository(db, log)
	id := []string{
		"1271d9a1-d540-4b10-a438-13c1276fc9e2",
		"add84ccd-1a69-4c53-97cc-484fda5302d6",
	}
	product, err := productRepository.FindProductsById(db, id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(product)
	// validate := config.NewValidator(viperConfig)
	// app := config.NewFiber(viperConfig)
	// config.Bootstrap(&config.BootstrapConfig{
	// 	DB:       db,
	// 	App:      app,
	// 	Log:      log,
	// 	Validate: validate,
	// 	Config:   viperConfig,
	// })
}
