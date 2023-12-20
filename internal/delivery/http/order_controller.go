package http

import (
	"order-service/internal/model"
	"order-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OrderController struct {
	Log     *logrus.Logger
	Usecase usecase.OrderUsecase
}

func NewOrderController(useCase usecase.OrderUsecase, logger *logrus.Logger) *OrderController {
	return &OrderController{
		Log:     logger,
		Usecase: useCase,
	}
}

func (c *OrderController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateOrderRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	response, err := c.Usecase.Insert(ctx.UserContext(), request)

	if err != nil {
		c.Log.WithError(err).Error("failed to create order")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.OrderResponse]{Data: response, Code: fiber.StatusOK, Status: "OK"})
}
