package mq

import (
	"encoding/json"
	"errors"
	"ticket-seckill/service"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ngaut/log"
)

var RedisMq *redisMq

type redisMq struct {
	orderSerivce service.IOrderService
	redis        *redis.Client
}

func (mq *redisMq) Send(goodsId, userId int64, orderId string) (err error) {
	msg := MqMsg{
		GoodsId: goodsId,
		UserId:  userId,
		OrderId: orderId,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if err = mq.redis.LPush(ctx, MqName, jsonData).Err(); err != nil {
		return err
	}

	return nil
}

func (mq *redisMq) Receive() {
	for {
		reply, err := mq.redis.RPop(ctx, MqName).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				// 列表为空，没有消息可接收
				time.Sleep(100 * time.Millisecond)
				continue
			}
			log.Errorf("redis.ZRangeByScore() failed, err: %v", err)
			continue
		}
		// 将接收到的 JSON 数据反序列化为 Message 结构体
		var msg MqMsg
		if err := json.Unmarshal([]byte(reply), &msg); err != nil {
			log.Errorf("Unmarshal failed, err: %v", err)
			continue
		}
		if err = mq.orderSerivce.CreateOrder(msg.GoodsId, msg.UserId, msg.OrderId); err != nil {
			log.Errorf("CreateOrder failed, err: %v", err)
			continue
		}

	}
}
