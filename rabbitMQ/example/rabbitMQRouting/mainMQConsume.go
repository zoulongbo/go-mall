package main

import (
	"fmt"
	"github.com/zoulongbo/go-mall/rabbitMQ"
)

func main()  {
	fmt.Println("routing consume 1")
	fmt.Println()
	rabbitMq := rabbitMQ.NewRabbitMQRouting("myRouting", "routing1")
	rabbitMq.ConsumeRouting()
}
