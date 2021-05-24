package main

import (
	"fmt"
	"github.com/zoulongbo/go-mall/rabbitMQ"
)

func main()  {
	fmt.Println("topic consume 2")
	fmt.Println()
	rabbitMq := rabbitMQ.NewRabbitMQTopic("myTopic", "#")
	rabbitMq.ConsumeTopic()
}
