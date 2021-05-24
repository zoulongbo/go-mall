package main

import (
	"github.com/zoulongbo/go-mall/rabbitMQ"
	"strconv"
	"time"
)

func main()  {
	rabbitMqOne := rabbitMQ.NewRabbitMQRouting("myRouting", "routing1")
	rabbitMqTwo := rabbitMQ.NewRabbitMQRouting("myRouting", "routing2")
	for i:=1; i < 100; i++ {
		rabbitMqOne.PublishRouting("i am routing1: " + strconv.Itoa(i)+"th")
		rabbitMqTwo.PublishRouting("i am routing2: " + strconv.Itoa(i)+"th")

		time.Sleep(time.Second)
	}
}

