package service

import "ticket-seckill/model"

var (
	OrderService IOrderService
	GoodsService IGoodsService
)

func GetOrderService() IOrderService {
	return OrderService
}

func GetGoodsService() IGoodsService {
	return GoodsService
}

type IOrderService interface {
	Seckill(int64, int64) error
	CreateOrderNavie(int64, int64) error
}

type IGoodsService interface {
	GetGoods(id int64) (model.Goods, error)
	SeckillNavie(userId, goodsId int64) error
}
