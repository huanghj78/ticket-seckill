package mq

import (
	"encoding/json"
	"log"

	"ticket-seckill/service"

	"github.com/streadway/amqp"
)

var RabbitMq *rabbitMq

type rabbitMq struct {
	orderService service.IOrderService
	conn         *amqp.Connection
	channel      *amqp.Channel
	queueName    string
}

func (mq *rabbitMq) Send(goodsId, userId int64, orderId string) error {
	msg := MqMsg{
		GoodsId: goodsId,
		UserId:  userId,
		OrderId: orderId,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = mq.channel.Publish(
		"",           // exchange
		mq.queueName, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonData,
		})
	if err != nil {
		return err
	}

	return nil
}

func (mq *rabbitMq) Receive() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	_, err = ch.QueueDeclare(
		mq.queueName, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := mq.channel.Consume(
		mq.queueName, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for d := range msgs {
		var msg MqMsg
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		if err := mq.orderService.CreateOrder(msg.GoodsId, msg.UserId, msg.OrderId); err != nil {
			log.Printf("CreateOrder failed: %v", err)
			continue
		}
	}
}
