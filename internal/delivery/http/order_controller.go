package http

import (
	"math"
	"order-service/internal/model"
	"order-service/internal/pkg"
	"order-service/internal/pkg/exception"
	"order-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OrderController struct {
	Log      *logrus.Logger
	Usecase  usecase.OrderUsecase
	Validate *pkg.Validation
}

func NewOrderController(useCase usecase.OrderUsecase, logger *logrus.Logger, validate *pkg.Validation) *OrderController {
	return &OrderController{
		Log:      logger,
		Usecase:  useCase,
		Validate: validate,
	}
}

func (c *OrderController) Search(ctx *fiber.Ctx) error {
	request := &model.SearchOrderRequest{
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}
	responses, total, _ := c.Usecase.Search(ctx.UserContext(), request)
	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}
	// if err.Status != exception.ERRBADREQUEST {
	// 	c.Log.WithError(err.Error).Error("failed to search order")
	// 	errValidation := c.Validate.ErrorJson(err.Error)
	// 	return ctx.JSON(model.ErrorResponse[pkg.ErrorMessage]{
	// 		Code:   int64(err.Status),
	// 		Status: fiber.ErrBadRequest.Message,
	// 		Error:  errValidation,
	// 	})
	// }

	return ctx.JSON(model.WebResponse[[]model.OrderResponse]{Data: responses, Paging: paging, Code: fiber.StatusOK, Status: "OK"})
}

func (c *OrderController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateOrdersRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return err
	}

	response, err := c.Usecase.Insert(ctx.UserContext(), request)

	if err.Status == exception.ERRBADREQUEST {
		c.Log.WithError(err.Error).Error("failed to create order")
		errValidation := c.Validate.ErrorJson(err.Error)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse[interface{}]{
			Code:   fiber.StatusBadRequest,
			Status: fiber.ErrBadRequest.Message,
			Error:  errValidation,
		})
		// arr := map[string][]string{
		// 	"email": {
		// 		"required",
		// 		"must email",
		// 	},
		// 	"password": {
		// 		"required",
		// 	},
		// }
		// return ctx.Status(fiber.StatusBadRequest).JSON(arr)
	}

	if err != nil {
		c.Log.WithError(err.Error).Error("failed to create order")
		return err.Error
	}

	return ctx.JSON(model.WebResponse[*model.OrderResponse]{Data: response, Code: fiber.StatusOK, Status: "OK"})
}
