package order

import (
	"sync"
	"ticket-seckill/repository"
	"ticket-seckill/service"

	"github.com/go-redis/redis/v8"
)

var once sync.Once

func InitService() {
	once.Do(func() {
		service.OrderService = &orderService{
			orderRepository: repository.NewOrderRepository(),
			goodsService:    service.GetGoodsService(),
		}
	})
}

type orderService struct {
	orderRepository repository.OrderRepository
	goodsService    service.IGoodsService
	redis           *redis.Client
}
