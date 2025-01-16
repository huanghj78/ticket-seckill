package model

import (
	"time"
)

const (
	CLOSED int8 = -1
	UNPAID int8 = 0
	PAYING int8 = 1
	PAID   int8 = 2
)

type Order struct {
	Id      uint64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	OrderId string `gorm:"type:varchar(32);uniqueIndex:idx_uid_gid;not null" json:"order_id"`
	UserId  int64  `gorm:"not null" json:"user_id"`
	GoodsId int64  `gorm:"not null" json:"goods_id"`
}

type OrderInfo struct {
	Id         int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	OrderId    string    `gorm:"type:varchar(32);uniqueIndex:idx_order_id;not null" json:"order_id"`
	UserId     int64     `gorm:"not null" json:"user_id"`
	GoodsId    int64     `gorm:"not null" json:"goods_id"`
	GoodsName  string    `gorm:"type:varchar(128);not null" json:"goods_name"`
	GoodsImg   string    `gorm:"type:varchar(128);not null" json:"goods_img"`
	GoodsPrice int64     `gorm:"not null" json:"goods_price"`
	Status     int8      `gorm:"default:0;not null" json:"status"`
	CreateTime time.Time `gorm:"autoCreateTime;not null" json:"create_time"`
	UpdateTime time.Time `gorm:"autoUpdateTime;not null" json:"update_time"`
}

func NewOrderInfo(userId int64, goods Goods, orderId string) OrderInfo {
	order := OrderInfo{
		OrderId:    orderId,
		UserId:     userId,
		GoodsId:    goods.Id,
		GoodsName:  goods.Name,
		GoodsImg:   goods.Img,
		GoodsPrice: goods.Price,
		Status:     UNPAID,
	}
	return order
}

type OrderCount struct {
	Unfinished int64 `db:"unfinished" json:"unfinished"`
	Finished   int64 `db:"finished" json:"finished"`
	Closed     int64 `db:"closed" json:"closed"`
}

type SeckillResult struct {
	Status  int8   `json:"status"`
	OrderId string `json:"orderId,omitempty"`
}
