package controller

import (
	"net/http"

	v1 "github.com/kolibriee/trade-metrics/internal/controller/http/v1"
	"github.com/kolibriee/trade-metrics/internal/repository"
)

type Controller struct {
	Handler http.Handler
}

func NewController(repo *repository.Repository) *Controller {
	return &Controller{
		Handler: v1.NewHandler(repo).InitRouterGin(),
	}
}
