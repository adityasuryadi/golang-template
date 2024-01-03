package usecase

import (
	"context"
	"errors"
	"fmt"
	"order-service/internal/entity"
	"order-service/internal/model"
	"order-service/internal/model/converter"
	"order-service/internal/pkg"
	"order-service/internal/pkg/exception"
	"order-service/internal/repository"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderUsecaseImpl struct {
	OrderRepository   repository.OrderRepository
	ProductRepository repository.Productrepository
	Log               *logrus.Logger
	DB                *gorm.DB
	Validate          *pkg.Validation
	JaegerTracer      *pkg.JaegerTracer
}

func NewOrderUsecase(db *gorm.DB, log *logrus.Logger, repository repository.OrderRepository, productRepo repository.Productrepository, tracer *pkg.JaegerTracer, validate *pkg.Validation) *OrderUsecaseImpl {
	return &OrderUsecaseImpl{
		OrderRepository:   repository,
		Log:               log,
		DB:                db,
		Validate:          validate,
		ProductRepository: productRepo,
		JaegerTracer:      tracer,
	}
}

// Search implements OrderUsecase.
func (u *OrderUsecaseImpl) Search(ctx context.Context, request *model.SearchOrderRequest) ([]model.OrderResponse, int64, *exception.CustomError) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := u.Validate.ValidateRequest(request); err != nil {
		// u.Log.WithError(err).Error("error validating request body")
		return nil, 0, &exception.CustomError{
			Status: exception.ERRBADREQUEST,
			Error:  err,
		}
	}

	orders, total, err := u.OrderRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Error("error getting contacts")
		return nil, 0, &exception.CustomError{
			Status: exception.ERRSERVER,
			Error:  err,
		}
	}

	responses := make([]model.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = *converter.OrderToResponse(&order)
	}
	return responses, total, nil
}

func findElement(slice []entity.Product, key interface{}) entity.Product {
	var product entity.Product
	for _, v := range slice {
		if v.Id == key {
			product = v
		}
	}
	return product
}

// Insert implements OrderUsecase.
func (u *OrderUsecaseImpl) Insert(ctx context.Context, request *model.CreateOrdersRequest) (*model.OrderResponse, *exception.CustomError) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	_, span := u.JaegerTracer.Tracer.Start(ctx, "create-order")
	defer span.End()

	err := u.Validate.ValidateRequest(request)
	if err != nil {
		return nil, &exception.CustomError{
			Status: exception.ERRBADREQUEST,
			Error:  err,
		}
	}

	span.AddEvent("Generate Bill No")
	orderNo, orderNoCounter := GenerateBillNo(tx)
	orderProducts := []entity.OrderProduct{}

	// get product by id's
	ids := []string{}
	for _, v := range request.Orders {
		ids = append(ids, v.ProductId)
	}

	span.AddEvent("find product")
	products, err := u.ProductRepository.FindProductsById(tx, ids)
	if err != nil {
		u.Log.WithError(err).Error("failed to create order")
		span.RecordError(err)
	}

	var totalOrderPrice float64
	for _, v := range request.Orders {
		product := findElement(products, uuid.MustParse(v.ProductId))
		orderProducts = append(orderProducts, entity.OrderProduct{
			ProductId:    uuid.MustParse(v.ProductId),
			Qty:          v.Qty,
			ProductName:  product.Name,
			ProductPrice: product.Price,
			TotalPrice:   float64(v.Qty) * product.Price,
		})
		totalOrderPrice += float64(v.Qty) * product.Price
	}

	order := &entity.Order{
		OrderNo:          orderNo,
		OrderNoCounter:   orderNoCounter,
		OrderProducts:    orderProducts,
		TotalPrice:       totalOrderPrice,
		ShippmentAddress: "Bandung",
	}

	span.AddEvent("Insert Order")
	err = u.OrderRepository.Create(tx, order)
	if err != nil {
		span.RecordError(err)
		u.Log.WithError(err).Error("failed to create order")
		return nil, &exception.CustomError{
			Status: exception.ERRSERVER,
			Error:  err,
		}
	}

	if err := tx.Commit().Error; err != nil {
		span.RecordError(err)
		u.Log.WithError(err).Error("failed to create order")
		return nil, &exception.CustomError{
			Status: exception.ERRSERVER,
			Error:  err,
		}
	}

	orderProductResponses := []model.OrderProductResponse{}
	for _, v := range order.OrderProducts {
		orderProductResponses = append(orderProductResponses, model.OrderProductResponse{
			Id:           v.Id.String(),
			ProductId:    v.ProductId.String(),
			ProductName:  v.ProductName,
			ProductPrice: strconv.Itoa(int(v.ProductPrice)),
			Qty:          v.Qty,
		})
	}

	return &model.OrderResponse{
		Id:            order.Id.String(),
		OrderNo:       order.OrderNo,
		OrderProducts: orderProductResponses,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	}, nil
}

// generate Order Number or invoice

func GenerateBillNo(tx *gorm.DB) (orderNo string, orderNoCounter int64) {
	order := new(entity.Order)
	today := time.Now().Format("2006-01-02")
	curdate := time.Now().Format("20060102")
	result := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(clause.Expr{SQL: "to_timestamp(created_at/1000)::date = ?", Vars: []interface{}{today}}).
		Order("order_no_counter desc").First(order)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		orderNo = fmt.Sprintf("INV-%s%d", curdate, 1)
		orderNoCounter = 1
	} else {
		orderNoCounter += order.OrderNoCounter + 1
		orderNo = fmt.Sprintf("INV-%s%d", curdate, orderNoCounter)
	}
	return orderNo, orderNoCounter
}
