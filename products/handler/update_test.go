package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_UpdateProduct(t *testing.T) {
	ctx := context.Background()
	store := setupTest(t)
	defer store.client.Close()

	userId := uuid.New()
	product := addProduct(t, ctx, store.client, userId)
	product2 := addProduct(t, ctx, store.client, uuid.New())
	productReserved := addProduct(t, ctx, store.client, userId)
	cookie := setupCookie(t, userId)

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "it returns a 404 if the provided id does not exist",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%s", uuid.New()), nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.NotEqual(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "it returns a 401 if the user is not authenticated",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%s", uuid.New()), nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "it returns a 401 if the user does not own the product",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"title": "Test Title",
					"price": 20,
				}
				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%s", product2.ID), bytes.NewBuffer(reqBody))
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
			name: "it returns a 400 if the product is reserved",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"title": "Test Title",
					"price": 20,
				}
				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				err = productReserved.Update().SetOrderID(uuid.New()).Exec(ctx)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%s", productReserved.ID), bytes.NewBuffer(reqBody))
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
			name: "it returns an error if an invalid title is provided",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"title": "",
					"price": 20,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)
				return http.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%s", product.ID), bytes.NewBuffer(reqBody))
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
			name: "it returns an error if an invalid price is provided",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"title": "Test title",
					"price": -10,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)
				return http.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(reqBody))
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
			name: "it updated the product with valid inputs",
			setupRequest: func() (*http.Request, error) {
				store.mockPublisher.On("PublishProductUpdated", mock.AnythingOfType("*ent.Product")).Return()

				data := gin.H{
					"title": "New Title",
					"price": 30,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)
				return http.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%s", product.ID), bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				store.mockPublisher.AssertCalled(t, "PublishProductUpdated", mock.AnythingOfType("*ent.Product"))
				respBody := &ProductResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, product.ID.String(), respBody.ID)
				assert.Equal(t, "New Title", respBody.Title)
				assert.Equal(t, 30, respBody.Price)
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
