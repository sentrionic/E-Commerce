package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	status "github.com/sentrionic/ecommerce/common/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v72"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreatePayment(t *testing.T) {
	store := setupTest(t)
	defer store.client.Close()
	ctx := context.Background()

	userId := uuid.New()
	order := addOrder(t, ctx, store.client, userId)
	order2 := addOrder(t, ctx, store.client, uuid.New())
	orderCancelled := addOrder(t, ctx, store.client, userId)
	err := orderCancelled.Update().SetStatus(status.Cancelled).Exec(ctx)
	assert.NoError(t, err)

	cookie := setupCookie(t, userId)

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "returns a 404 when purchasing an order that does not exist",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"token":   "testtoken",
					"orderId": uuid.New(),
				}
				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)
				return http.NewRequest(http.MethodPost, "/api/payments", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "returns a 401 when purchasing an order that does not belong to the user",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"token":   "testtoken",
					"orderId": order2.ID,
				}
				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)
				return http.NewRequest(http.MethodPost, "/api/payments", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "returns a 400 when purchasing a cancelled order",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"token":   "testtoken",
					"orderId": orderCancelled.ID,
				}
				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)
				return http.NewRequest(http.MethodPost, "/api/payments", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "it creates a product with valid inputs",
			setupRequest: func() (*http.Request, error) {
				store.mockService.On("HandleCharge", order.Price, "testtoken").Return(&stripe.Charge{ID: "id"}, nil)
				store.mockPublisher.On("PublishPaymentCreated", mock.AnythingOfType("*ent.Payment")).Return()

				data := gin.H{
					"token":   "testtoken",
					"orderId": order.ID,
				}
				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)
				return http.NewRequest(http.MethodPost, "/api/payments", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)
				store.mockPublisher.AssertExpectations(t)
				store.mockService.AssertExpectations(t)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := tc.setupRequest()
			assert.NoError(t, err)
			tc.setupHeaders(t, request)
			store.router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}
