package model

import (
	"ticket-seckill/infra/code"
	"time"
)

const (
	ENDED      int8 = -1 // 已结束
	NOTSTARTED int8 = 0  // 未开始
	ONGOING    int8 = 1  // 进行中
	SOLDOUT    int8 = 2  // 已售罄
)

type Goods struct {
	Id          int64     `gorm:"primaryKey;autoIncrement;comment:'商品id'"`   // `id` 字段，主键且自增
	Name        string    `gorm:"type:varchar(255);not null;comment:'商品名称'"` // `name` 字段，非空
	Img         string    `gorm:"type:varchar(255);not null;comment:'商品图片'"` // `img` 字段，非空
	OriginPrice int64     `gorm:"not null;comment:'商品价格'"`                   // `origin_price` 字段，非空
	Price       int64     `gorm:"not null;comment:'秒杀价格'"`                   // `price` 字段，非空
	Stock       int       `gorm:"not null;comment:'库存'"`                     // `stock` 字段，非负
	StartTime   time.Time `gorm:"not null;comment:'秒杀开始时间'"`                 // `start_time` 字段，非空
	EndTime     time.Time `gorm:"not null;comment:'秒杀结束时间'"`
}

type GoodsVO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Img         string `json:"img"`
	OriginPrice int64  `json:"originPrice"`
	Price       int64  `json:"price"`
	Duration    int64  `json:"duration"`
	Status      int8   `json:"status"`
}

func (goods Goods) ToVO() GoodsVO {
	g := GoodsVO{}
	g.Id = goods.Id
	g.OriginPrice = goods.OriginPrice
	g.Name = goods.Name
	g.Img = goods.Img
	g.Price = goods.Price
	startTime := goods.StartTime.Unix()
	endTime := goods.EndTime.Unix()
	now := time.Now().Unix()
	if now < startTime {
		g.Status = NOTSTARTED
		g.Duration = startTime - now
	} else if now >= startTime && now <= endTime {
		if goods.Stock > 0 {
			g.Status = ONGOING
		} else {
			g.Status = SOLDOUT
		}
	} else {
		g.Status = ENDED
	}
	return g
}

func (goods Goods) Check() (err error) {
	now := time.Now().Unix()
	startTime := goods.StartTime.Unix()
	endTime := goods.EndTime.Unix()
	if now < startTime {
		err = code.MiaoshaNotStart
	} else if now > endTime {
		err = code.MiaoshaEnded
	} else if goods.Stock <= 0 {
		err = code.GoodsSaleOut
	}
	return
}
