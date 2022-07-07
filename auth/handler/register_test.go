package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/ecommerce/auth/ent"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Register(t *testing.T) {

	ctx := context.Background()
	store := setupTest(t)
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {

		}
	}(store.client)

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "It returns a 201 on successful register",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": "password",
					"username": "test user",
					"email":    "test@test.com",
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "It returns a 400 with an invalid email",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": "password",
					"username": "test user",
					"email":    "tester",
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "It returns a 400 with an invalid password",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": "p",
					"username": "test user",
					"email":    "test@example.com",
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "It returns a 400 with missing email",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": "password",
					"username": "test user",
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "It returns a 400 with missing password",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"email":    "test@example.com",
					"username": "test user",
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "It disallows duplicate emails",
			setupRequest: func() (*http.Request, error) {
				_, err := store.client.User.Create().
					SetEmail("test@example.com").
					SetPassword("password").
					SetUsername("username").
					Save(ctx)
				assert.NoError(t, err)

				data := gin.H{
					"email":    "test@example.com",
					"password": "password",
					"username": "test user",
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "It sets a cookie after successful register",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": "password",
					"username": "test user",
					"email":    "test@test.com",
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)
				assert.Contains(t, recorder.Header(), "Set-Cookie")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			_, err := store.client.User.Delete().Exec(ctx)
			assert.NoError(t, err)
			request, err := tc.setupRequest()
			assert.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			store.router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}
