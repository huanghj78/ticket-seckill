package handler

import (
	"net/http"
	"ticket-seckill/service"

	"github.com/gin-gonic/gin"
	"github.com/ngaut/log"
)

type OrderHandler struct {
	orderService service.IOrderService
}

func InitOrderHandler() *OrderHandler {
	return &OrderHandler{orderService: service.GetOrderService()}
}

func (h *OrderHandler) Seckill(c *gin.Context) {
	r := new(struct {
		UserId  int64 `json:"userId" binding:"required"`
		GoodsId int64 `json:"goodsId" binding:"required"`
	})
	if err := c.Bind(r); err != nil {
		log.Error(err)
		return
	}
	if err := h.orderService.Seckill(r.UserId, r.GoodsId); err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, nil)
	}

}
