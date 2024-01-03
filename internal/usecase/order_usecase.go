package usecase

import (
	"context"
	"order-service/internal/model"
	"order-service/internal/pkg/exception"
)

type OrderUsecase interface {
	Insert(ctx context.Context, request *model.CreateOrdersRequest) (*model.OrderResponse, *exception.CustomError)
	Search(ctx context.Context, request *model.SearchOrderRequest) ([]model.OrderResponse, int64, *exception.CustomError)
}
