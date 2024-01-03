package converter

import (
	"order-service/internal/entity"
	"order-service/internal/model"
)

func OrderToResponse(order *entity.Order) *model.OrderResponse {
	orderProductResponses := []model.OrderProductResponse{}
	for _, v := range order.OrderProducts {
		orderProductResponses = append(orderProductResponses, *OrderProductToResponse(&v))
	}
	return &model.OrderResponse{
		Id:            order.Id.String(),
		OrderNo:       order.OrderNo,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
		OrderProducts: orderProductResponses,
	}
}
