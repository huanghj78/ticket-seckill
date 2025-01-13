package router

import (
	"ticket-seckill/handler"
	"ticket-seckill/service/goods"
	"ticket-seckill/service/order"

	"github.com/gin-gonic/gin"
)

var (
	goodsHandler *handler.GoodsHandler
	orderHandler *handler.OrderHandler
)

func initService() {
	goods.InitService()
	order.InitService()
}

func initHandler() {
	goodsHandler = handler.InitGoodsHandler()
	orderHandler = handler.InitOrderHandler()
}

func initRouter(router *gin.Engine) {
	router.POST("/seckill", orderHandler.Seckill)
}

func Init() (router *gin.Engine) {
	initService()
	initHandler()
	router = gin.Default()
	initRouter(router)
	return
}
