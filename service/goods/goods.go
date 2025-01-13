package goods

import (
	"sync"
	"ticket-seckill/infra/cache"
	"ticket-seckill/repository"
	"ticket-seckill/service"

	"github.com/go-redis/redis/v8"
)

var once sync.Once

func InitService() {
	once.Do(func() {
		service.GoodsService = &goodsService{
			goodsRepository: repository.NewGoodsRepository(),
			redis:           cache.Client,
		}
	})
}

type goodsService struct {
	goodsRepository repository.GoodsRepository
	redis           *redis.Client
}
