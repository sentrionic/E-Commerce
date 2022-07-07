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

func TestHandler_Login(t *testing.T) {

	ctx := context.Background()
	store := setupTest(t)
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {

		}
	}(store.client)

	_, err := store.client.User.Create().
		SetEmail("test@test.com").
		SetPassword("password").
		SetUsername("username").
		Save(ctx)
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "It responds with a cookie when given valid credentials",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": "password",
					"email":    "test@test.com",
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(reqBody))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Contains(t, recorder.Header(), "Set-Cookie")
			},
		},
		{
			name: "It returns a 404 with a non existing user",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": "password",
					"email":    "tester@example.com",
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(reqBody))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "It returns a 404 with an invalid password",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": "password123",
					"email":    "test@test.com",
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(reqBody))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
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
