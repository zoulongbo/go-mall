package services

import (
	"github.com/zoulongbo/go-mall/models"
	"github.com/zoulongbo/go-mall/repositories"
)

type Order interface {
	GetOrderById(id int64) (order *models.Order, err error)
	GetAllOrder() (order []*models.Order, err error)
	InsertOrder(order *models.Order) (id int64, err error)
	DeleteOrderById(id int64) bool
	UpdateOrder(order *models.Order) error
	GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderByMessage(*models.Message) (int64, error)
}

type OrderService struct {
	orderRepos repositories.Order
}

func NewOrderService() Order {
	return &OrderService{orderRepos: repositories.NewOrderManager()}
}

func (o *OrderService) GetOrderById(id int64) (order *models.Order, err error) {
	return o.orderRepos.SelectByKey(id)
}

func (o *OrderService) GetAllOrder() (order []*models.Order, err error) {
	return o.orderRepos.SelectAll()
}

func (o *OrderService) InsertOrder(order *models.Order) (id int64, err error) {
	return o.orderRepos.Insert(order)
}

func (o *OrderService) DeleteOrderById(id int64) bool {
	return o.orderRepos.Delete(id)
}

func (o *OrderService) UpdateOrder(order *models.Order) error {
	return o.orderRepos.Update(order)
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.orderRepos.SelectAllWithInfo()
}

func (o *OrderService) InsertOrderByMessage(message *models.Message) (orderId int64, err error) {
	order := &models.Order{
		UserId:      message.UserId,
		ProductId:   message.ProductId,
		OrderStatus: models.OrderStatusSuccess,
	}
	return o.InsertOrder(order)
}
