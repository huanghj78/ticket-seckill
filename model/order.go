package model

import "time"

const (
	CLOSED int8 = -1
	UNPAID int8 = 0
	PAYING int8 = 1
	PAID   int8 = 2
)

type Order struct {
	Id      int64  `db:"id"`
	OrderId string `db:"order_id"`
	UserId  int64  `db:"user_id"`
	GoodsId int64  `db:"goods_id"`
}

type OrderInfo struct {
	Id         int64     `db:"id"`
	OrderId    string    `db:"order_id"`
	UserId     int64     `db:"user_id"`
	GoodsId    int64     `db:"goods_id"`
	GoodsName  string    `db:"goods_name"`
	GoodsImg   string    `db:"goods_img"`
	GoodsPrice int64     `db:"goods_price"`
	PaymentId  int64     `db:"payment_id"`
	Status     int8      `db:"status"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

func NewOrderInfo(userId int64, goods Goods) OrderInfo {
	order := OrderInfo{
		OrderId:    createOrderId(),
		UserId:     userId,
		GoodsId:    goods.Id,
		GoodsName:  goods.Name,
		GoodsImg:   goods.Img,
		GoodsPrice: goods.Price,
		Status:     UNPAID,
	}
	return order
}

func createOrderId() string {
	return time.Now().Format("20060102150405")
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
