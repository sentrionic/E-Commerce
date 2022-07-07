package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetOrder(t *testing.T) {
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
				return http.NewRequest(http.MethodGet, fmt.Sprintf("/api/orders/%s", order.ID), nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "it returns the order",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, fmt.Sprintf("/api/orders/%s", order.ID), nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				respBody := &OrderResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.Equal(t, order.ID.String(), respBody.ID)
			},
		},
		{
			name: "returns an error if one user tries to fetch another users order",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, fmt.Sprintf("/api/orders/%s", order2.ID), nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
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
