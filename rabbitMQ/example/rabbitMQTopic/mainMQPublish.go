package main

import (
	"github.com/zoulongbo/go-mall/rabbitMQ"
	"strconv"
	"time"
)

func main()  {
	rabbitMqOne := rabbitMQ.NewRabbitMQTopic("myTopic", "x.am")
	rabbitMqTwo := rabbitMQ.NewRabbitMQTopic("myTopic", "y.ou")
	for i:=1; i < 100; i++ {
		rabbitMqOne.PublishTopic("i am topic1: " + strconv.Itoa(i)+"th")
		rabbitMqTwo.PublishTopic("i am topic2: " + strconv.Itoa(i)+"th")

		time.Sleep(time.Second)
	}
}

