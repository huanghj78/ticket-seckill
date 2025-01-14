package repository

import (
	"ticket-seckill/infra/db"
	"ticket-seckill/model"

	"github.com/ngaut/log"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderNavie(model.OrderInfo) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() OrderRepository {
	return &orderRepository{db: db.DB}
}

func (r *orderRepository) CreateOrderNavie(orderInfo model.OrderInfo) (err error) {
	txn := r.db.Begin()
	defer func() {
		if err != nil {
			txn.Rollback()
		} else {
			txn.Commit()
		}
	}()
	// 创建订单
	log.Info(orderInfo)
	var order model.Order
	order.OrderId = orderInfo.OrderId
	order.GoodsId = orderInfo.GoodsId
	order.UserId = orderInfo.UserId
	log.Infof("%v\n", order)
	if err = r.db.Save(&order).Error; err != nil {
		return
	}

	if err = r.db.Save(&orderInfo).Error; err != nil {
		return
	}
	return

}
