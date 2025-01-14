package repository

import (
	"errors"
	"log"
	"ticket-seckill/infra/db"
	"ticket-seckill/model"

	"gorm.io/gorm"
)

type GoodsRepository interface {
	GetGoods(int64) (model.Goods, error)
	SeckillNaive(int64, int64) error
}

type goodsRepository struct {
	db *gorm.DB
}

func NewGoodsRepository() GoodsRepository {
	return &goodsRepository{db: db.DB}
}

func (r *goodsRepository) GetGoods(id int64) (goods model.Goods, err error) {
	if err = r.db.Where("id = ?", id).First(&goods).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil // 如果没有找到记录，返回 nil 错误
		} else {
			log.Printf("r.db.Get() failed, err: %v", err)
		}
	}
	return
}

func (r *goodsRepository) SeckillNaive(userId, goodsId int64) (err error) {
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

	var goods model.Goods
	if err = txn.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", goodsId).First(&goods).Error; err != nil {
		return err
	}

	if err = goods.Check(); err != nil {
		return err
	}

	if goods.Stock <= 0 {
		return errors.New("stock not sufficient")
	}

	goods.Stock -= 1
	if err = txn.Save(&goods).Error; err != nil {
		return err
	}

	return nil
}
