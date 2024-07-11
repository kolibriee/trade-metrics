package v1

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kolibriee/trade-metrics/internal/domain"
	"github.com/kolibriee/trade-metrics/internal/repository"
	mock_repository "github.com/kolibriee/trade-metrics/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_GetOrderHistory(t *testing.T) {
	type mockBehavior func(r *mock_repository.Mockorderhistory, client *domain.Client)

	tests := []struct {
		name                 string
		queryParams          string
		inputClient          *domain.Client
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			queryParams: "client-name=Misha&exchange-name=binance&label=111&pair=BTCUSDT",
			inputClient: &domain.Client{
				ClientName:   "Misha",
				ExchangeName: "binance",
				Label:        "111",
				Pair:         "BTCUSDT",
			},
			mockBehavior: func(r *mock_repository.Mockorderhistory, client *domain.Client) {
				r.EXPECT().GetOrderHistory(client).Return([]*domain.HistoryOrder{
					{
						Client:              *client,
						Side:                "buy",
						Type:                "limit",
						BaseQty:             1.0,
						Price:               50000.0,
						AlgorithmNamePlaced: "alg1",
						LowestSellPrice:     49900.0,
						HighestBuyPrice:     50100.0,
						CommissionQuoteQty:  0.1,
						TimePlaced:          time.Time{},
					},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"client":{"client_name":"Misha","exchange_name":"binance","label":"111","pair":"BTCUSDT"},"side":"buy","type":"limit","base_qty":1,"price":50000,"algorithm_name_placed":"alg1","lowest_sell_prc":49900,"highest_buy_prc":50100,"commission_quote_qty":0.1,"time_placed":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name:                 "Invalid Input",
			queryParams:          "client-name=&exchange-name=&label=&pair=",
			inputClient:          &domain.Client{},
			mockBehavior:         func(r *mock_repository.Mockorderhistory, client *domain.Client) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input"}`,
		},
		{
			name:        "Server Error",
			queryParams: "client-name=Misha&exchange-name=binance&label=111&pair=BTCUSDT",
			inputClient: &domain.Client{
				ClientName:   "Misha",
				ExchangeName: "binance",
				Label:        "111",
				Pair:         "BTCUSDT",
			},
			mockBehavior: func(r *mock_repository.Mockorderhistory, client *domain.Client) {
				r.EXPECT().GetOrderHistory(client).Return(nil, errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockorderhistory(c)
			tt.mockBehavior(repo, tt.inputClient)

			handler := NewHandler(&repository.Repository{Orderhistory: repo})

			r := gin.New()
			r.GET("/orderhistory", handler.GetOrderHistory)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/orderhistory?"+tt.queryParams, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_SaveOrder(t *testing.T) {
	type mockBehavior func(r *mock_repository.Mockorderhistory, order *domain.HistoryOrder)

	tests := []struct {
		name                 string
		inputBody            string
		inputOrder           *domain.HistoryOrder
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputBody: `{
				"client": {"client_name": "Misha", "exchange_name": "binance", "label": "test", "pair": "BTCUSDT"},
				"side": "buy", "type": "limit", "base_qty": 1.0, "price": 50000.0, "algorithm_name_placed": "alg1",
				"lowest_sell_prc": 49900.0, "highest_buy_prc": 50100.0, "commission_quote_qty": 0.1
			}`,
			inputOrder: &domain.HistoryOrder{
				Client: domain.Client{
					ClientName:   "Misha",
					ExchangeName: "binance",
					Label:        "test",
					Pair:         "BTCUSDT",
				},
				Side:                "buy",
				Type:                "limit",
				BaseQty:             1.0,
				Price:               50000.0,
				AlgorithmNamePlaced: "alg1",
				LowestSellPrice:     49900.0,
				HighestBuyPrice:     50100.0,
				CommissionQuoteQty:  0.1,
				TimePlaced:          time.Time{},
			},
			mockBehavior: func(r *mock_repository.Mockorderhistory, order *domain.HistoryOrder) {
				r.EXPECT().SaveOrder(gomock.Any()).Return(nil).Do(func(order *domain.HistoryOrder) {
					order.TimePlaced = time.Time{}
				})
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:                 "Invalid Input Body",
			inputBody:            `{"client": {"client_name": "", "exchangeame": "binance", "label": "test", "pair": "BTCUSDT"}}`,
			inputOrder:           &domain.HistoryOrder{},
			mockBehavior:         func(r *mock_repository.Mockorderhistory, order *domain.HistoryOrder) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name: "Server Error",
			inputBody: `{
				"client": {"client_name": "Misha", "exchange_name": "binance", "label": "test", "pair": "BTCUSDT"},
				"side": "buy", "type": "limit", "base_qty": 1.0, "price": 50000.0, "algorithm_name_placed": "alg1",
				"lowest_sell_prc": 49900.0, "highest_buy_prc": 50100.0, "commission_quote_qty": 0.1
			}`,
			inputOrder: &domain.HistoryOrder{
				Client: domain.Client{
					ClientName:   "Misha",
					ExchangeName: "binance",
					Label:        "test",
					Pair:         "BTCUSDT",
				},
				Side:                "buy",
				Type:                "limit",
				BaseQty:             1.0,
				Price:               50000.0,
				AlgorithmNamePlaced: "alg1",
				LowestSellPrice:     49900.0,
				HighestBuyPrice:     50100.0,
				CommissionQuoteQty:  0.1,
				TimePlaced:          time.Now(), // This will be set dynamically in the handler
			},
			mockBehavior: func(r *mock_repository.Mockorderhistory, order *domain.HistoryOrder) {
				r.EXPECT().SaveOrder(gomock.Any()).Return(errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockorderhistory(c)
			tt.mockBehavior(repo, tt.inputOrder)

			handler := NewHandler(&repository.Repository{Orderhistory: repo})

			r := gin.New()
			r.POST("/order", handler.SaveOrder)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/order", bytes.NewBufferString(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
