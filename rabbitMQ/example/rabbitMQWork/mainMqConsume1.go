package main

import (
	"fmt"
	"github.com/zoulongbo/go-mall/rabbitMQ"
)

func main()  {
	fmt.Println("i am consume1")
	rabbitMq := rabbitMQ.NewRabbitMQSimple("myFirst")
	rabbitMq.ConsumeSimple()
}