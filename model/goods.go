package model

import "time"

const (
	ENDED      int8 = -1 // 已结束
	NOTSTARTED int8 = 0  // 未开始
	ONGOING    int8 = 1  // 进行中
	SOLDOUT    int8 = 2  // 已售罄
)

type Goods struct {
	Id          int64     `db:"id"`
	Name        string    `db:"name"`
	Img         string    `db:"img"`
	OriginPrice int64     `db:"origin_price"`
	Price       int64     `db:"price"`
	Stock       int       `db:"stock"`
	StartTime   time.Time `db:"start_time"`
	EndTime     time.Time `db:"end_time"`
}
