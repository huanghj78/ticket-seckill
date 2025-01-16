package order

import (
	"fmt"
	"strconv"
	"sync"
	"ticket-seckill/infra/cache"
	"ticket-seckill/infra/code"
	"ticket-seckill/model"
	"ticket-seckill/mq"
	"ticket-seckill/repository"
	"ticket-seckill/service"
	"time"

	"github.com/ngaut/log"

	"github.com/go-redis/redis/v8"
)

var once sync.Once
var mu sync.Mutex

func InitService() {
	once.Do(func() {
		service.OrderService = &orderService{
			orderRepository: repository.NewOrderRepository(),
			goodsService:    service.GetGoodsService(),
			redis:           cache.Client,
		}
	})
}

type orderService struct {
	orderRepository repository.OrderRepository
	goodsService    service.IGoodsService
	redis           *redis.Client
}

func (s *orderService) Seckill(userId, goodsId int64) (err error) {
	var (
		goods   model.Goods
		orderId string
	)
	if goods, err = s.goodsService.GetGoods(goodsId); err != nil {
		return
	}
	// 校验秒杀开始、结束时间
	if err = goods.Check(); err != nil {
		return
	}
	// 加锁 确保一人一单
	if err = s.tryLock(userId, goodsId); err != nil {
		return
	}

	if err = s.decrStock(goodsId, userId); err != nil {
		return
	}

	// 生成订单号
	if orderId, err = s.generateOrderId(goodsId); err != nil {
		s.incrStock(goodsId, userId)
		return
	}

	mq.RedisMq.Send(goodsId, userId, orderId)

	if err = s.UnLock(userId, goodsId); err != nil {
		return
	}

	return

}

func (s *orderService) decrStock(goodsId, userId int64) (err error) {
	luaScript := `
	local stock = redis.call('GET', KEYS[1])
	if tonumber(stock) <= 0 then
		return 1
	end
	local uIdExists = redis.call('SISMEMBER', KEYS[2], ARGV[1])
	if uIdExists == 1 then
		return 2
	end
	redis.call('DECR', KEYS[1])
	redis.call('SADD', KEYS[2], ARGV[1])
	return 0
	`

	script := redis.NewScript(luaScript)
	res, err := script.Run(service.Ctx, s.redis, []string{"stock:" + strconv.FormatInt(goodsId, 10), "orders:" + strconv.FormatInt(goodsId, 10)}, strconv.FormatInt(userId, 10)).Int()
	if err != nil {
		return
	}
	if res == 1 {
		log.Infof("%d 库存不足", goodsId)
		err = fmt.Errorf("库存不足")
		return
	} else if res == 2 {
		log.Infof("%d %d用户重复下单", goodsId, userId)
		err = fmt.Errorf("用户重复下单")
		return
	}
	return
}

func (s *orderService) incrStock(goodsId, userId int64) (err error) {
	luaScript := `
	redis.call('IECR', KEYS[1])
	redis.call('SREM', KEYS[2], ARGV[1])
	return 0
	`
	script := redis.NewScript(luaScript)
	_, err = script.Run(service.Ctx, s.redis, []string{"stock:" + strconv.FormatInt(goodsId, 10), "orders:" + strconv.FormatInt(goodsId, 10)}, strconv.FormatInt(userId, 10)).Int()
	return
}

// func (s *orderService) Seckill(userId, goodsId int64) (err error) {
// 	mu.Lock()         // 加锁
// 	defer mu.Unlock() // 解锁
// 	if err = s.goodsService.SeckillNavie(userId, goodsId); err != nil {
// 		return
// 	}
// 	err = s.CreateOrderNavie(userId, goodsId)
// 	return
// }

// func (s *orderService) CreateOrderNavie(userId, goodsId int64) (err error) {
// 	var goods model.Goods
// 	// 查询商品信息
// 	if goods, err = s.goodsService.GetGoods(goodsId); err != nil {
// 		log.Errorf("goodsService.GetGoods() failed, err: %v", err)
// 		return
// 	}

// 	// 创建订单
// 	orderInfo := model.NewOrderInfo(userId, goods)
// 	if err = s.orderRepository.CreateOrderNavie(orderInfo); err != nil {
// 		log.Errorf("orderRepository.CreateOrderNavie() failed, err: %v", err)
// 		return
// 	}
// 	return

// }

func (s *orderService) tryLock(userId, goodsId int64) (err error) {
	lockId := 1
	var res bool
	if res, err = s.redis.SetNX(service.Ctx, fmt.Sprintf("lock:%d:%d", userId, goodsId), lockId, time.Minute).Result(); err != nil {
		log.Errorf("redis.SetNx() failed, err: %v", err)
		err = code.RedisErr
		return
	}
	if !res {
		err = code.MiaoshaFailed
	}
	return
}

func (s *orderService) UnLock(userId, goodsId int64) (err error) {
	if err = s.redis.Del(service.Ctx, fmt.Sprintf("lock:%d:%d", userId, goodsId)).Err(); err != nil {
		log.Errorf("redis.Del() failed, err: %v", err)
		err = code.RedisErr
	}
	return
}

func (s *orderService) GetOrderId(userId, goodsId int64) (orderId string, err error) {
	if orderId, err = s.redis.Get(service.Ctx, fmt.Sprintf(service.OrderUidGidKey, userId, goodsId)).Result(); err != nil {
		if err == redis.Nil {
			err = nil
		} else {
			log.Errorf("redis.Get() failed, err: %v", err)
			err = code.RedisErr
		}
	}
	return
}

func (s *orderService) GetOrderInfo(orderId string) (orderInfo model.OrderInfo, err error) {
	if orderInfo, err = s.orderRepository.GetOrderInfo(orderId); err != nil {
		err = code.DBErr
	}
	return
}

func (s *orderService) generateOrderId(goodsId int64) (string, error) {
	// 获取当前时间的秒级时间戳
	timestamp := strconv.Itoa(int(time.Now().Unix()))

	// 构建Redis键名，包含商品ID和时间戳
	key := fmt.Sprintf("icr:%d:%s", goodsId, time.Now().Format("2006-01-02"))

	// 使用Redis的INCR命令递增计数器
	seq, err := s.redis.Incr(service.Ctx, key).Result()
	if err != nil {
		return "", err // 如果Redis操作失败，返回错误
	}

	// 构建订单ID，格式可以是：时间戳_序列号
	orderId := fmt.Sprintf("%s%d", timestamp, seq)
	return orderId, nil
}

func (s *orderService) CreateOrder(goodsId, userId int64, orderId string) (err error) {
	if err = s.orderRepository.CreateOrder(goodsId, userId, orderId); err != nil {
		log.Errorf("orderRepository.CreateOrder() failed, err: %v", err)
		return
	}
	return
}
