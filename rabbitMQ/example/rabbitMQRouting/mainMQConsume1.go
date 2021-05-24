package main

import (
	"fmt"
	"github.com/zoulongbo/go-mall/rabbitMQ"
)

func main()  {
	fmt.Println("routing consume 2")
	fmt.Println()
	rabbitMq := rabbitMQ.NewRabbitMQRouting("myRouting", "routing2")
	rabbitMq.ConsumeRouting()
}
