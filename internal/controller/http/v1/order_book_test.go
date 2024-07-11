package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kolibriee/trade-metrics/internal/domain"
	"github.com/kolibriee/trade-metrics/internal/repository"
	mock_repository "github.com/kolibriee/trade-metrics/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_GetOrderBook(t *testing.T) {
	type mockBehavior func(r *mock_repository.Mockorderbook, exchangeName, pair string)

	tests := []struct {
		name                 string
		exchange_name        string
		pair                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "OK",
			exchange_name: "binance",
			pair:          "BTCUSDT",
			mockBehavior: func(r *mock_repository.Mockorderbook, exchangeName, pair string) {
				r.EXPECT().GetOrderBook(exchangeName, pair).Return(&domain.AsksBids{
					Id: 0,
					Asks: []domain.DepthOrder{
						{
							Price:   50000,
							BaseQty: 1,
						},
					},
					Bids: []domain.DepthOrder{
						{
							Price:   51000,
							BaseQty: 2,
						},
					},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":0,"asks":[{"price":50000,"base_qty":1}],"bids":[{"price":51000,"base_qty":2}]}`,
		},
		{
			name:                 "empty input",
			exchange_name:        "binance",
			pair:                 "",
			mockBehavior:         func(r *mock_repository.Mockorderbook, exchangeName, pair string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input"}`,
		},
		{
			name:          "server error",
			exchange_name: "binance",
			pair:          "BTCUSDT",
			mockBehavior: func(r *mock_repository.Mockorderbook, exchangeName, pair string) {
				r.EXPECT().GetOrderBook(exchangeName, pair).Return(nil, errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mock_repository.NewMockorderbook(c)
			tt.mockBehavior(repo, tt.exchange_name, tt.pair)
			handler := NewHandler(&repository.Repository{Orderbook: repo})
			r := gin.New()
			r.GET("/orderbook/:exchangeName/:pair/", handler.GetOrderBook)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/orderbook/%s/%s/", tt.exchange_name, tt.pair), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_SaveOrderBook(t *testing.T) {
	type mockBehavior func(r *mock_repository.Mockorderbook, exchangeName, pair string, asksBids *domain.AsksBids)

	tests := []struct {
		name                 string
		exchange_name        string
		pair                 string
		inputBody            string
		inputAsksBids        *domain.AsksBids
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "OK",
			exchange_name: "binance",
			pair:          "BTCETH",

			inputBody: `{"asks":[{"price":100,"base_qty":1},{"price":110,"base_qty":5}],"bids":[{"price":99,"base_qty":2}]}`,
			inputAsksBids: &domain.AsksBids{
				Id:   12345,
				Asks: []domain.DepthOrder{{Price: 100, BaseQty: 1}, {Price: 110, BaseQty: 5}},
				Bids: []domain.DepthOrder{{Price: 99, BaseQty: 2}},
			},
			mockBehavior: func(r *mock_repository.Mockorderbook, exchangeName, pair string, asksBids *domain.AsksBids) {
				r.EXPECT().SaveOrderBook(exchangeName, pair, gomock.Any()).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":"12345"}`,
		},
		{
			name:          "Invalid input",
			exchange_name: "",
			pair:          "BTCETH",
			inputBody:     `{"asks":[{"price":"100","base_qty":"1"}],"bids":[]`,
			inputAsksBids: &domain.AsksBids{
				Asks: []domain.DepthOrder{{Price: 100, BaseQty: 1}},
			},
			mockBehavior:         func(r *mock_repository.Mockorderbook, exchangeName, pair string, asksBids *domain.AsksBids) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input"}`,
		},
		{
			name:                 "Invalid input body",
			exchange_name:        "binance",
			pair:                 "BTCETH",
			inputBody:            `{"asks":[],"bids":[]`,
			inputAsksBids:        &domain.AsksBids{},
			mockBehavior:         func(r *mock_repository.Mockorderbook, exchangeName, pair string, asksBids *domain.AsksBids) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:          "Server Error",
			exchange_name: "binance",
			pair:          "BTCETH",
			inputBody:     `{"asks":[{"price":100,"base_qty":1},{"price":110,"base_qty":5}],"bids":[{"price":99,"base_qty":2}]}`,
			inputAsksBids: &domain.AsksBids{
				Asks: []domain.DepthOrder{{Price: 100, BaseQty: 1}, {Price: 110, BaseQty: 5}},
				Bids: []domain.DepthOrder{{Price: 99, BaseQty: 2}},
			},
			mockBehavior: func(r *mock_repository.Mockorderbook, exchangeName, pair string, asksBids *domain.AsksBids) {
				r.EXPECT().SaveOrderBook(exchangeName, pair, gomock.Any()).Return(errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mock_repository.NewMockorderbook(c)
			tt.mockBehavior(repo, tt.exchange_name, tt.pair, tt.inputAsksBids)
			handler := NewHandler(&repository.Repository{Orderbook: repo})
			r := gin.New()
			r.POST("/orderbook/:exchangeName/:pair/", handler.SaveOrderBook)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/orderbook/%s/%s/", tt.exchange_name, tt.pair), bytes.NewBufferString(tt.inputBody))
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if tt.name == "OK" {
				assert.Equal(t, tt.expectedResponseBody, `{"id":"12345"}`)
			} else {
				assert.Equal(t, tt.expectedResponseBody, w.Body.String())
			}
		})
	}
}
