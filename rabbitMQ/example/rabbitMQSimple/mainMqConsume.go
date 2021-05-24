package main

import "github.com/zoulongbo/go-mall/rabbitMQ"

func main()  {
	rabbitMq := rabbitMQ.NewRabbitMQSimple("myFirst")
	rabbitMq.ConsumeSimple()
}