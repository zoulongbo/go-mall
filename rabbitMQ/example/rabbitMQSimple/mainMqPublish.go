package main

import (
	"github.com/zoulongbo/go-mall/rabbitMQ"
)

func main()  {
	rabbitMQ := rabbitMQ.NewRabbitMQSimple("myFirst")
	rabbitMQ.PublishSimple("hello world")
}

