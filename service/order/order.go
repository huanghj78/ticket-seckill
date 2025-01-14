package order

import (
	"log"
	"sync"
	"ticket-seckill/model"
	"ticket-seckill/repository"
	"ticket-seckill/service"

	"github.com/go-redis/redis/v8"
)

var once sync.Once
var mu sync.Mutex

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

func (s *orderService) Seckill(userId, goodsId int64) (err error) {
	mu.Lock()         // 加锁
	defer mu.Unlock() // 解锁
	if err = s.goodsService.SeckillNavie(userId, goodsId); err != nil {
		return
	}
	err = s.CreateOrderNavie(userId, goodsId)
	return
}

func (s *orderService) CreateOrderNavie(userId, goodsId int64) (err error) {
	var goods model.Goods
	// 查询商品信息
	if goods, err = s.goodsService.GetGoods(goodsId); err != nil {
		log.Printf("goodsService.GetGoods() failed, err: %v", err)
		return
	}

	// 创建订单
	orderInfo := model.NewOrderInfo(userId, goods)
	if err = s.orderRepository.CreateOrderNavie(orderInfo); err != nil {
		log.Printf("orderRepository.CreateOrderNavie() failed, err: %v", err)
		return
	}
	return

}
