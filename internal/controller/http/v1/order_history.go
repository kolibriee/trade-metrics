package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kolibriee/trade-metrics/internal/domain"
)

func (h *Handler) GetOrderHistory(c *gin.Context) {
	clientName := c.Query("client-name")
	exchangeName := c.Query("exchange-name")
	label := c.Query("label")
	pair := c.Query("pair")
	if clientName == "" || exchangeName == "" || label == "" || pair == "" {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input").Error())
		return
	}

	client := domain.Client{
		ClientName:   clientName,
		ExchangeName: exchangeName,
		Label:        label,
		Pair:         pair,
	}
	orderHistory, err := h.repo.GetOrderHistory(&client)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("server error").Error())
		return
	}
	c.JSON(http.StatusOK, orderHistory)
}

func (h *Handler) SaveOrder(c *gin.Context) {
	var order domain.HistoryOrder
	if err := c.BindJSON(&order); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input body").Error())
		return
	}
	order.TimePlaced = time.Now()
	if err := h.repo.SaveOrder(&order); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("server error").Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
