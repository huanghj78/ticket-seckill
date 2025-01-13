package service

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
}

type IGoodsService interface {
}
