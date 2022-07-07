package handler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetOrders(t *testing.T) {
	ctx := context.Background()
	store := setupTest(t)
	defer store.client.Close()

	userId := uuid.New()
	userId2 := uuid.New()
	cookie := setupCookie(t, userId)

	// Add one order for user2
	product1 := addProduct(t, ctx, store.client)
	addOrder(t, ctx, store.client, product1, userId2)

	// Add two orders for user1
	product2 := addProduct(t, ctx, store.client)
	product3 := addProduct(t, ctx, store.client)
	order2 := addOrder(t, ctx, store.client, product2, userId)
	order3 := addOrder(t, ctx, store.client, product3, userId)

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "It can fetch a list of orders",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/orders", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]OrderResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				r := *respBody

				assert.Equal(t, 2, len(r))
				assert.Equal(t, order2.ID.String(), r[0].ID)
				assert.Equal(t, order3.ID.String(), r[1].ID)
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
