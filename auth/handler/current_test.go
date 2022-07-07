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

func TestHandler_CurrentUser(t *testing.T) {

	ctx := context.Background()
	store := setupTest(t)
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {

		}
	}(store.client)

	email := "test@example.com"
	password := "password"

	_, err := store.client.User.Create().
		SetEmail(email).
		SetPassword(password).
		SetUsername("username").
		Save(ctx)
	assert.NoError(t, err)

	cookie := ""

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Setup: Login User",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": password,
					"email":    email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(reqBody))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Contains(t, recorder.Header(), "Set-Cookie")
				cookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "It responds with the details about the current user",
			setupRequest: func() (*http.Request, error) {
				router, err := http.NewRequest(http.MethodGet, "/api/auth/current", nil)
				assert.NoError(t, err)
				router.Header.Add("Cookie", cookie)
				return router, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "It responds with an unauthorized error",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/auth/current", nil)
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
			request.Header.Set("Content-Type", "application/json")
			store.router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}
