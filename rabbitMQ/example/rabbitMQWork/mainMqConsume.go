package main

import (
	"fmt"
	"github.com/zoulongbo/go-mall/rabbitMQ"
)

func main()  {
	fmt.Println("i am consume")
	rabbitMq := rabbitMQ.NewRabbitMQSimple("myFirst")
	rabbitMq.ConsumeSimple()
}