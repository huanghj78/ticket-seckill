package mq

import (
	"context"
	"log"
	"ticket-seckill/service"

	"github.com/streadway/amqp"
)

var ctx = context.Background()

const MqName = "message_queue"

type MqMsg struct {
	GoodsId int64  `json:"goods_id"`
	UserId  int64  `json:"user_id"`
	OrderId string `json:"order_id"`
}

type IMessageQueue interface {
	Send(int64, int64, string) error
	Receive()
}

func Init() {
	// RedisMq = &redisMq{
	// 	orderSerivce: service.GetOrderService(),
	// 	redis:        cache.Client,
	// }
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	queueName := "seckill"

	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	RabbitMq = &rabbitMq{
		orderService: service.GetOrderService(),
		conn:         conn,
		channel:      ch,
		queueName:    queueName,
	}
}

func Run() {
	// go RedisMq.Receive()

	go RabbitMq.Receive()
}
