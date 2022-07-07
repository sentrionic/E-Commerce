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

func TestHandler_GetProduct(t *testing.T) {
	ctx := context.Background()
	store := setupTest(t)
	defer store.client.Close()

	userId := uuid.New()
	product := addProduct(t, ctx, store.client, userId)

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "it returns a 404 if the product is not found",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, fmt.Sprintf("/api/products/%s", uuid.New()), nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "it returns the product if the product is found",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, fmt.Sprintf("/api/products/%s", product.ID), nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				respBody := &ProductResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, product.ID.String(), respBody.ID)
				assert.Equal(t, product.Title, respBody.Title)
				assert.Equal(t, product.Price, respBody.Price)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := tc.setupRequest()
			assert.NoError(t, err)
			store.router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}
