package model

type OrderResponse struct {
	Id            string                 `json:"id"`
	OrderNo       string                 `json:"order_no"`
	CreatedAt     int64                  `json:"created_at"`
	UpdatedAt     int64                  `json:"updated_at"`
	OrderProducts []OrderProductResponse `json:"order_product"`
}

type OrderProductResponse struct {
	Id           string `json:"id"`
	ProductId    string `json:"product_id"`
	ProductName  string `json:"product_name"`
	ProductPrice string `json:"product_price"`
	Qty          int64  `json:"qty"`
}

type CreateOrdersRequest struct {
	Orders []CreateOrderRequest `json:"orders" validate:"required,dive"`
}

type CreateOrderRequest struct {
	ProductId string `json:"product_id" validate:"required"`
	Qty       int64  `json:"qty" validate:"required"`
}

type SearchOrderRequest struct {
	Page int `json:"page" validate:"required,min=1"`
	Size int `json:"size" validate:"required,min=1,max=100"`
}
