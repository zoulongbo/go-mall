package main

import (
	"github.com/zoulongbo/go-mall/rabbitMQ"
	"strconv"
)

func main()  {
	rabbitMq := rabbitMQ.NewRabbitMQSub("mySub")
	for i:=0; i < 200; i++ {
		rabbitMq.PublishSub("hello everyone: " + strconv.Itoa(i))
	}
}

