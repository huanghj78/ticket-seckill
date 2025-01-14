package goods

import (
	"sync"
	"ticket-seckill/infra/cache"
	"ticket-seckill/infra/code"
	"ticket-seckill/model"
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

func (s *goodsService) GetGoods(id int64) (goods model.Goods, err error) {
	if goods, err = s.goodsRepository.GetGoods(id); err != nil {
		err = code.DBErr
		return
	}
	return
}

func (s *goodsService) SeckillNavie(usdeId, goodsId int64) (err error) {
	err = s.goodsRepository.SeckillNaive(usdeId, goodsId)
	return
}
