package usecase

import (
	"context"
	"fmt"
	"order-service/internal/entity"
	"order-service/internal/model"
	"order-service/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderUsecaseImpl struct {
	OrderRepository repository.OrderRepository
	Log             *logrus.Logger
	DB              *gorm.DB
}

// Insert implements OrderUsecase.
func (u *OrderUsecaseImpl) Insert(ctx context.Context, request *model.CreateOrderRequest) (*model.OrderResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	fmt.Println(ctx)
	order := &entity.Order{
		ProductId:    uuid.MustParse(request.ProductId),
		ProductName:  "susus",
		ProductPrice: 2500,
		Qty:          request.Qty,
	}
	err := u.OrderRepository.Create(tx, order)
	if err != nil {
		u.Log.WithError(err).Error("failed to create order")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("failed to create order")
		return nil, fiber.ErrInternalServerError
	}

	return &model.OrderResponse{
		Id:        order.Id.String(),
		ProductId: order.ProductId.String(),
		Qty:       order.Qty,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}, nil
}

func NewOrderUsecase(db *gorm.DB, log *logrus.Logger, repository repository.OrderRepository) OrderUsecase {
	return &OrderUsecaseImpl{
		OrderRepository: repository,
		Log:             log,
		DB:              db,
	}
}
