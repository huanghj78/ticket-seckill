package repository

import (
	"errors"
	"fmt"
	"ticket-seckill/infra/db"
	"ticket-seckill/model"

	"github.com/ngaut/log"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderNavie(model.OrderInfo) error
	GetOrderInfo(string) (model.OrderInfo, error)
	CloseOrder(model.OrderInfo) error
	CreateOrder(goodsId, userId int64, orderId string) (err error)
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
		if p := recover(); p != nil {
			txn.Rollback()
			panic(p)
		} else if err != nil {
			txn.Rollback()
		} else if commitErr := txn.Commit().Error; commitErr != nil {
			err = fmt.Errorf("commit failed: %w", commitErr)
		}
	}()

	// 创建订单
	log.Info(orderInfo)
	var order model.Order
	order.OrderId = orderInfo.OrderId
	order.GoodsId = orderInfo.GoodsId
	order.UserId = orderInfo.UserId

	if err = txn.Create(&order).Error; err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	if err = txn.Create(&orderInfo).Error; err != nil {
		return fmt.Errorf("failed to create order info: %w", err)
	}
	return nil
}

func (r *orderRepository) GetOrderInfo(orderId string) (order model.OrderInfo, err error) {
	// 查询订单
	result := r.db.First(&order, "order_id = ?", orderId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return order, errors.New("order not found")
		}
		return order, result.Error
	}
	return

}

func (r *orderRepository) CloseOrder(order model.OrderInfo) (err error) {
	txn := r.db.Begin()
	defer func() {
		if p := recover(); p != nil {
			txn.Rollback()
			panic(p)
		} else if err != nil {
			txn.Rollback()
		} else if commitErr := txn.Commit().Error; commitErr != nil {
			err = fmt.Errorf("commit failed: %w", commitErr)
		}
	}()

	// 加库存
	if err = txn.Model(&model.Goods{}).
		Where("id = ?", order.GoodsId).
		UpdateColumn("stock", gorm.Expr("stock + ?", 1)).Error; err != nil {
		return err
	}

	// 删除订单
	if err = txn.Where("order_id = ?", order.OrderId).
		Delete(&model.Order{}).Error; err != nil {
		return err
	}

	// 修改订单信息状态
	if err = txn.Model(&model.OrderInfo{}).
		Where("order_id = ? AND status = ?", order.OrderId, model.UNPAID).
		Update("status", model.CLOSED).Error; err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) CreateOrder(goodsId, userId int64, orderId string) (err error) {
	txn := r.db.Begin()
	defer func() {
		if p := recover(); p != nil {
			txn.Rollback()
			panic(p)
		} else if err != nil {
			txn.Rollback()
		} else {
			err = txn.Commit().Error
		}
	}()
	// 减库存
	var goods model.Goods
	if err := txn.First(&goods, goodsId).Error; err != nil {
		return err
	}
	if goods.Stock <= 0 {
		log.Error("超卖！！！")
		return errors.New("insufficient stock")
	}
	goods.Stock -= 1
	if err := txn.Save(&goods).Error; err != nil {
		return err
	}
	// 创建订单
	var order model.Order
	order.GoodsId = goodsId
	order.UserId = userId
	order.OrderId = orderId
	if err := txn.Create(&order).Error; err != nil {
		return err
	}

	// 创建订单信息
	orderInfo := model.NewOrderInfo(userId, goods, orderId)
	if err := txn.Create(&orderInfo).Error; err != nil {
		return err
	}
	return
}
