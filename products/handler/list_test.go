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

func TestHandler_GetProducts(t *testing.T) {
	ctx := context.Background()
	store := setupTest(t)
	defer store.client.Close()

	for i := 0; i < 3; i++ {
		addProduct(t, ctx, store.client, uuid.New())
	}

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "It can fetch a list of products",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/products", nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]ProductResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, 3, len(*respBody))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := tc.setupRequest()
			assert.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			store.router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}
