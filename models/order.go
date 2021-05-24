package models

type Order struct {
	ID          int64 `json:"id" sql:"id" info:"id"`
	UserId      int64 `json:"userId" sql:"user_id"  info:"userId"`
	ProductId   int64 `json:"productId" sql:"product_id" info:"productId"`
	OrderStatus int   `json:"orderStatus" sql:"order_status" info:"orderStatus"`
}

const (
	OrderTable = "order"

	OrderStatusWait    = iota
	OrderStatusSuccess //不赋值默认 iot+1
	OrderStatusFailed  //不赋值默认 iot+1+1
)
