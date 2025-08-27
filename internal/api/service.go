package api

import (
	"github.com/biryanim/wb_tech_L0/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type Implementation struct {
	orderService service.OrderService
	//consumerService service.ConsumerService
}

func NewImplementation(orderService service.OrderService) *Implementation {
	return &Implementation{
		orderService: orderService,
	}
}

func (i *Implementation) GetOrder(c *gin.Context) {
	orderUIDStr := c.Param("order_uid")
	orderUID, err := uuid.Parse(orderUIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order UUID format",
		})
		return
	}

	order, err := i.orderService.GetOrder(c.Request.Context(), orderUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, order)
}
