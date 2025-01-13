package repository

import (
	"ticket-seckill/infra/db"

	"gorm.io/gorm"
)

type OrderRepository interface {
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() OrderRepository {
	return &orderRepository{db: db.DB}
}
