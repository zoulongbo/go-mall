package main

import (
	"fmt"
	"github.com/zoulongbo/go-mall/rabbitMQ"
)

func main()  {
	fmt.Println("topic consume 1")
	fmt.Println()
	rabbitMq := rabbitMQ.NewRabbitMQTopic("myTopic", "x.*")
	rabbitMq.ConsumeTopic()
}
