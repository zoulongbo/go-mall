package main

import (
	"github.com/zoulongbo/go-mall/rabbitMQ"
	"strconv"
)

func main()  {
	rabbitMQ := rabbitMQ.NewRabbitMQSimple("myFirst")
	for i:=0; i < 200; i++ {
		rabbitMQ.PublishSimple("hello world: " + strconv.Itoa(i))
	}
}

