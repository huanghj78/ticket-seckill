package mq

import (
	"context"
	"ticket-seckill/infra/cache"
	"ticket-seckill/service"
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
	RedisMq = &redisMq{
		orderSerivce: service.GetOrderService(),
		redis:        cache.Client,
	}
}

func Run() {
	go RedisMq.Receive()
}
