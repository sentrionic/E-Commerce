package handler

import (
	"github.com/sentrionic/ecommerce/auth/ent"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Logout(t *testing.T) {

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
			name: "It clears the cookie after logout",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Contains(t, recorder.Header(), "Set-Cookie")
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
