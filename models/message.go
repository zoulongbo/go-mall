package models

type Message struct {
	UserId    int64	`json:"userId"`
	ProductId int64	`json:"productId"`
}

func NewMessage(userId, productId int64) *Message {
	return &Message{
		UserId:    userId,
		ProductId: productId,
	}
}
