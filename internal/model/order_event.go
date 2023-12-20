package model

type OrderEvent struct {
	Id           string `json:"id"`
	OrderNo      string `json:"order_no"`
	ProductId    string `json:"product_id"`
	ProductName  string `json:"product_name"`
	ProductPrice string `json:"product_price"`
	Qty          int64  `json:"qty"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}
