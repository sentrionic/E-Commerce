package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	status "github.com/sentrionic/ecommerce/common/order"
	gen "github.com/sentrionic/ecommerce/orders/ent/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_DeleteOrder(t *testing.T) {
	ctx := context.Background()
	store := setupTest(t)
	defer store.client.Close()

	userId := uuid.New()
	product := addProduct(t, ctx, store.client)
	product2 := addProduct(t, ctx, store.client)
	order := addOrder(t, ctx, store.client, product, userId)
	order2 := addOrder(t, ctx, store.client, product2, uuid.New())
	cookie := setupCookie(t, userId)

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "it returns 401 if the user is not authenticated",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/orders/%s", order.ID), nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "it returns 401 if the user is not the owner of the order",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/orders/%s", order2.ID), nil)
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
			name: "it marks an order as cancelled",
			setupRequest: func() (*http.Request, error) {
				store.mockPublisher.On("PublishOrderCancelled", mock.AnythingOfType("*ent.Order")).Return()
				return http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/orders/%s", order.ID), nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				o, err := store.client.Order.Query().Where(gen.IDEQ(order.ID)).First(ctx)
				assert.NoError(t, err)
				assert.Equal(t, o.Status, status.Cancelled)
				store.mockPublisher.AssertCalled(t, "PublishOrderCancelled", mock.AnythingOfType("*ent.Order"))
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
