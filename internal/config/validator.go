package config

import (
	"order-service/internal/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *pkg.Validation {
	return pkg.NewValidation(validator.New())
}
