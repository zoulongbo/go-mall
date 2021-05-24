package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/zoulongbo/go-mall/models"
	"github.com/zoulongbo/go-mall/rabbitMQ"
	"github.com/zoulongbo/go-mall/services"
	"log"
)

func Order()  {
	//创建product serivce
	productService:=services.NewProductService()
	//创建order Service
	orderService := services.NewOrderService()
	//初始化 mq
	rabbitMQ := rabbitMQ.NewRabbitMQSimple("orderAdd")

	orderConsume(orderService,productService, rabbitMQ)
}

//消费 开始下单业务
func orderConsume(orderService services.Order,productService services.Product, rabbitMQ *rabbitMQ.RabbitMQ)  {
	//申请队列，不存在则创建，存在则创建
	_, err := rabbitMQ.Channel.QueueDeclare(rabbitMQ.QueueName, false, false,false, false, nil)
	if err != nil {
		fmt.Println(err)
	}
	rabbitMQ.Channel.Qos(
		//每次队列只消费一个消息 这个消息处理不完服务器不会发送第二个消息过来
		//当前消费者一次能接受的最大消息数量
		1,
		//服务器传递的最大容量
		0,
		//如果为true 对channel可用 false则只对当前队列可用
		false,)
	//处理完消息 手动应答
	msgs, err := rabbitMQ.Channel.Consume(rabbitMQ.QueueName, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			message := &models.Message{}
			log.Printf("Received a message: %s", d.Body)
			fmt.Println()
			//开始下单业务 数据转换 将[]byte转换为struct
			if err := json.Unmarshal(d.Body, message); err != nil {
				log.Printf("orderAdd json.Unmarshal failed : body %v, err:%v  ", d.Body, err)
				continue
			}
			product, err := productService.GetProductById(message.ProductId)
			if err != nil {
				log.Printf("orderAdd product get failed : message %v, err:%v  ", message, err)
				continue
			}
			//有库存再走下单
			if product.ProductNum > 0 {
				//缺锁 以防超卖
				err := productService.SubNumOne(message.ProductId, 1)
				if err != nil {
					log.Printf("orderAdd product num update failed : message %v, err:%v  ", message, err)
					continue
				}
				//创建订单
				orderId, err := orderService.InsertOrderByMessage(message)
				if err != nil {
					log.Printf("orderAdd  failed : message %v, err:%v  ", message, err)
					continue
				}
				log.Printf("orderAdd success orderId:%d", orderId)
				//false确认当前消息， true 确认所有消息
				d.Ack(false)
			}

		}
	}()

	log.Printf("[*] Wating for message, To exit press `CTRL + C`")
	<-forever
}