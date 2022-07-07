package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateProduct(t *testing.T) {
	store := setupTest(t)
	defer store.client.Close()

	cookie := setupCookie(t, uuid.New())

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "it has a route handler listening for post requests",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodPost, "/api/products", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.NotEqual(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "it can only be accessed if the user is signed in",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodPost, "/api/products", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
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
			name: "it creates a product with valid inputs",
			setupRequest: func() (*http.Request, error) {
				store.mockPublisher.On("PublishProductCreated", mock.AnythingOfType("*ent.Product")).Return()

				data := gin.H{
					"title": "Test Title",
					"price": 20,
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
				assert.Equal(t, http.StatusCreated, recorder.Code)
				store.mockPublisher.AssertCalled(t, "PublishProductCreated", mock.AnythingOfType("*ent.Product"))
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
