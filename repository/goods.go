package repository

import (
	"ticket-seckill/infra/db"

	"gorm.io/gorm"
)

type GoodsRepository interface {
}

type goodsRepository struct {
	db *gorm.DB
}

func NewGoodsRepository() GoodsRepository {
	return &goodsRepository{db: db.DB}
}
