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
	// CreateOrderNavie(int64, int64) error
	GetOrderInfo(string) (model.OrderInfo, error)
	CreateOrder(int64, int64, string) error
}

type IGoodsService interface {
	InitGoodsStock() error
	GetGoods(id int64) (model.Goods, error)
	SeckillNavie(userId, goodsId int64) error
	SetGoodsStock(goodsId int64, stock int) (err error)
	DecrStock(goodsId int64) (stock int64, err error)
	IncrStock(goodsId int64) (err error)
}
