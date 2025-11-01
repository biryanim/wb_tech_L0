package api

import (
	"net/http"

	"github.com/biryanim/wb_tech_L0/internal/service"
	"github.com/gin-gonic/gin"
)

// Implementation handles API endpoints for order operations.
type Implementation struct {
	orderService service.OrderService
}

// NewImplementation creates and returns a new Implementation instance with the provided OrderService.
func NewImplementation(orderService service.OrderService) *Implementation {
	return &Implementation{
		orderService: orderService,
	}
}

// GetOrder retrieves an order by its unique identifier from the request parameter and returns it as JSON.
func (i *Implementation) GetOrder(c *gin.Context) {
	orderUID := c.Param("order_uid")

	order, err := i.orderService.GetOrder(c.Request.Context(), orderUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, order)
}
