package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRouterGin() http.Handler {
	router := gin.New()
	router.Use(gin.Logger())
	orderBook := router.Group("/orderbook")
	{
		orderBook.GET("/:exchangeName/:pair/", h.GetOrderBook)
		orderBook.POST("/:exchangeName/:pair/", h.SaveOrderBook)
	}

	orderHistory := router.Group("/orderhistory")
	{
		orderHistory.GET("/", h.GetOrderHistory)
		orderHistory.POST("/", h.SaveOrder)
	}
	return router
}
