package converter

import (
	"order-service/internal/entity"
	"order-service/internal/model"
	"strconv"
)

func OrderProductToResponse(orderProduct *entity.OrderProduct) *model.OrderProductResponse {
	return &model.OrderProductResponse{
		Id:           orderProduct.Id.String(),
		ProductId:    orderProduct.ProductId.String(),
		ProductName:  orderProduct.ProductName,
		ProductPrice: strconv.Itoa(int(orderProduct.ProductPrice)),
		Qty:          orderProduct.Qty,
	}
}
