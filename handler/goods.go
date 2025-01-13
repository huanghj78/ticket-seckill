package handler

import "ticket-seckill/service"

type GoodsHandler struct {
	goodsService service.IGoodsService
}

func InitGoodsHandler() *GoodsHandler {
	return &GoodsHandler{goodsService: service.GetGoodsService()}
}
