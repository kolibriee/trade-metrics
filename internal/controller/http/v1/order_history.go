package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kolibriee/trade-metrics/internal/domain"
)

func (h *Handler) GetOrderHistory(c *gin.Context) {
	var client domain.Client
	if err := c.BindJSON(&client); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid json").Error()+err.Error())
		return
	}
	orderHistory, err := h.repo.GetOrderHistory(&client)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("failed to get order history").Error()+err.Error())
		return
	}
	c.JSON(http.StatusOK, orderHistory)
}

func (h *Handler) SaveOrder(c *gin.Context) {
	var order domain.HistoryOrder
	if err := c.BindJSON(&order); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid json").Error()+err.Error())
		return
	}
	order.TimePlaced = time.Now()
	if err := h.repo.SaveOrder(&order); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("failed to save order").Error()+err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
