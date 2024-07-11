package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kolibriee/trade-metrics/internal/domain"
)

func (h *Handler) GetOrderBook(c *gin.Context) {
	exchange := c.Param("exchangeName")
	pair := c.Param("pair")
	if exchange == "" || pair == "" {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input").Error())
		return
	}
	orderBook, err := h.repo.GetOrderBook(exchange, pair)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("server error").Error()+err.Error())
		return
	}
	c.JSON(http.StatusOK, orderBook)
}

func (h *Handler) SaveOrderBook(c *gin.Context) {
	exchange := c.Param("exchangeName")
	pair := c.Param("pair")
	if exchange == "" || pair == "" {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input").Error())
		return
	}
	var orderBook domain.AsksBids
	if err := c.BindJSON(&orderBook); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input body").Error())
		return
	}
	id := uuid.New().ID()
	orderBook.Id = id
	if err := h.repo.SaveOrderBook(exchange, pair, &orderBook); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, errors.New("server error").Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}
