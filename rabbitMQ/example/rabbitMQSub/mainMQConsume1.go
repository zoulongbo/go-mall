package main

import (
	"fmt"
	"github.com/zoulongbo/go-mall/rabbitMQ"
)

func main()  {
	fmt.Println("sub consume 1")
	fmt.Println()
	rabbitMq := rabbitMQ.NewRabbitMQSub("mySub")
	rabbitMq.ConsumeSub()
}
