package goods

import (
	"fmt"
	"log"
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
		service.GoodsService.InitGoodsStock()
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

func (s *goodsService) InitGoodsStock() (err error) {
	var list []model.Goods
	if list, err = s.goodsRepository.GetGoodsList(-1); err != nil {
		err = code.DBErr
		return
	}
	for _, v := range list {
		if err = s.SetGoodsStock(v.Id, v.Stock); err != nil {
			return
		}
	}
	return
}

func (s *goodsService) SetGoodsStock(goodsId int64, stock int) (err error) {
	if err = s.redis.Set(service.Ctx, fmt.Sprintf(service.GoodsStockKey, goodsId), stock, -1).Err(); err != nil {
		log.Printf("redis.Set() failed, err: %v", err)
		err = code.RedisErr
	}
	return
}

func (s *goodsService) DecrStock(goodsId int64) (stock int64, err error) {
	if stock, err = s.redis.Decr(service.Ctx, fmt.Sprintf(service.GoodsStockKey, goodsId)).Result(); err != nil {
		log.Printf("redis.Decr() failed, err: %v", err)
		err = code.RedisErr
	}
	return
}

func (s *goodsService) IncrStock(goodsId int64) (err error) {
	if err = s.redis.Incr(service.Ctx, fmt.Sprintf(service.GoodsStockKey, goodsId)).Err(); err != nil {
		log.Printf("redis.Incr() failed, err: %v", err)
		err = code.RedisErr
	}
	return
}
