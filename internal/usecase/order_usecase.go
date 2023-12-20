package usecase

import (
	"context"
	"order-service/internal/model"
)

type OrderUsecase interface {
	Insert(ctx context.Context, request *model.CreateOrderRequest) (*model.OrderResponse, error)
}
