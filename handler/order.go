package handler

import (
	"ticket-seckill/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.IOrderService
}

func InitOrderHandler() *OrderHandler {
	return &OrderHandler{orderService: service.GetOrderService()}
}

func (h *OrderHandler) Seckill(c *gin.Context) {
	// todo
}
