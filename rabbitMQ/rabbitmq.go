package rabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

//url格式 amqp://账号:密码@服务器地址:端口号/vhost
const MQURL = "amqp://user:123456@127.0.0.1:5672/test"

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel

	//队列名
	QueueName string
	//交换机
	Exchange string
	//key
	Key string
	//连接信息
	MqUrl string
	//锁
	sync.Mutex
}

//创建rabbitMQ实例
func NewRabbitMQ(queueName, exChange, key string) *RabbitMQ {
	rabbitMQ := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exChange,
		Key:       key,
		MqUrl:     MQURL,
	}
	var err error
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.MqUrl)
	rabbitMQ.failOnErr(err, "创建连接错误")

	rabbitMQ.Channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnErr(err, "获取Channel失败")
	return rabbitMQ
}

//断开rabbitMQ
func (r *RabbitMQ) Destroy() {
	r.Channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//simple模式 队列 + 消费者平均分配
//simple模式 1、创建简单实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}
//simple模式 2、生产者
func (r *RabbitMQ) PublishSimple(message string) error  {
	//申请队列，不存在则创建，存在则创建
	r.Lock()
	defer r.Unlock()

	_, err := r.Channel.QueueDeclare(r.QueueName, false, false,false, false, nil)
	if err != nil {
		return err
	}
	//发送消息到队列
	err = r.Channel.Publish(r.Exchange, r.QueueName, false, false, amqp.Publishing{
		ContentType:     "text/plain",
		Body:            []byte(message),
	})
	if err != nil {
		return err
	}
	return nil
}

//simple模式 3、消费者
func (r *RabbitMQ) ConsumeSimple()  {
	//申请队列，不存在则创建，存在则创建
	_, err := r.Channel.QueueDeclare(r.QueueName, false, false,false, false, nil)
	if err != nil {
		fmt.Println(err)
	}

	msgs, err := r.Channel.Consume(r.QueueName, "", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			fmt.Println()
			fmt.Printf("%s", d.Body)
		}
	}()

	log.Printf("[*] Wating for message, To exit press `CTRL + C`")
	<-forever
}
//	                       - c1
//发布订阅模式 一推多 p1-2-｜-  c2
//                         - c3
//发布订阅模式 1、创建发布订阅实例
func NewRabbitMQSub(exChange string) *RabbitMQ  {
	return NewRabbitMQ("", exChange, "")
}

//发布订阅模式 2、生产者
func (r *RabbitMQ) PublishSub(message string)  {
	//设置交换机  kind网络类型 fanout=广播
	err := r.Channel.ExchangeDeclare(r.Exchange, amqp.ExchangeFanout, false,false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange nge")
	//发送消息
	err = r.Channel.Publish(r.Exchange, "", false, false, amqp.Publishing{
		ContentType:     "text/plain",
		Body:            []byte(message),
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("publish success")
	fmt.Println()
}

//发布订阅模式 3、消费者
func (r *RabbitMQ) ConsumeSub()  {
	//交换机
	err := r.Channel.ExchangeDeclare(r.Exchange, amqp.ExchangeFanout, false,false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange")

	//随机生成队列
	q, err := r.Channel.QueueDeclare("", false, false,true, false, nil)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到exchange
	err = r.Channel.QueueBind(q.Name, "", r.Exchange, false, nil)
	r.failOnErr(err, "Failed to bind queue")

	//消费
	msgs, err := r.Channel.Consume(q.Name, "", true, false, false, false, nil)
	r.failOnErr(err, "Failed to consume")


	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			fmt.Println()
			fmt.Printf("%s", d.Body)
			fmt.Println()
		}
	}()

	log.Printf("sub model [*] Wating for message, To exit press `CTRL + C`")
	fmt.Println()
	<-forever
}


//routing模式 1、创建routing实例 可指定消费者
func NewRabbitMQRouting(exChange, routingKey string) *RabbitMQ  {
	return NewRabbitMQ("", exChange, routingKey)
}

//routing模式 2、生产者
func (r *RabbitMQ) PublishRouting(message string)  {
	//设置交换机  kind网络类型 fanout=广播
	err := r.Channel.ExchangeDeclare(r.Exchange, amqp.ExchangeDirect, false,false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange nge")
	//发送消息
	err = r.Channel.Publish(r.Exchange, r.Key, false, false, amqp.Publishing{
		ContentType:     "text/plain",
		Body:            []byte(message),
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("publish success")
	fmt.Println()
}

//routing模式 3、消费者
func (r *RabbitMQ) ConsumeRouting()  {
	//交换机
	err := r.Channel.ExchangeDeclare(r.Exchange, amqp.ExchangeDirect, false,false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange")

	//随机生成队列
	q, err := r.Channel.QueueDeclare("", false, false,true, false, nil)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到exchange
	err = r.Channel.QueueBind(q.Name, r.Key, r.Exchange, false, nil)
	r.failOnErr(err, "Failed to bind queue")

	//消费
	msgs, err := r.Channel.Consume(q.Name, "", true, false, false, false, nil)
	r.failOnErr(err, "Failed to consume")


	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			fmt.Println()
			fmt.Printf("%s", d.Body)
			fmt.Println()
		}
	}()

	log.Printf("Sub model [*] Wating for message, To exit press `CTRL + C`")
	fmt.Println()
	<-forever
}


//topic模式 1、创建topic实例 消费者可根据参数动态获取不同消息
func NewRabbitMQTopic(exChange, routingKey string) *RabbitMQ  {
	return NewRabbitMQ("", exChange, routingKey)
}

//Topic模式 2、生产者
func (r *RabbitMQ) PublishTopic(message string)  {
	//设置交换机  kind网络类型 fanout=广播
	err := r.Channel.ExchangeDeclare(r.Exchange, amqp.ExchangeTopic, false,false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange nge")
	//发送消息
	err = r.Channel.Publish(r.Exchange, r.Key, false, false, amqp.Publishing{
		ContentType:     "text/plain",
		Body:            []byte(message),
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("publish success")
	fmt.Println()
}

//Topic模式 3、消费者
func (r *RabbitMQ) ConsumeTopic()  {
	//交换机
	err := r.Channel.ExchangeDeclare(r.Exchange, amqp.ExchangeTopic, false,false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange")

	//随机生成队列
	q, err := r.Channel.QueueDeclare("", false, false,true, false, nil)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到exchange
	err = r.Channel.QueueBind(q.Name, r.Key, r.Exchange, false, nil)
	r.failOnErr(err, "Failed to bind queue")

	//消费
	msgs, err := r.Channel.Consume(q.Name, "", true, false, false, false, nil)
	r.failOnErr(err, "Failed to consume")


	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			fmt.Println()
			fmt.Printf("%s", d.Body)
			fmt.Println()
		}
	}()

	log.Printf("Topic model [*] Wating for message, To exit press `CTRL + C`")
	fmt.Println()
	<-forever
}
