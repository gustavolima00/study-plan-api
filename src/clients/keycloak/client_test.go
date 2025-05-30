package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"go-api/src/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

// MockHTTPClient is a mock implementation of http.Client
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetOIDCToken(t *testing.T) {
	tests := map[string]struct {
		HttpResponse       any
		HttpResponseStatus int
		ExpectedError      error
		Request            GetOIDCTokenRequest
		ExpectedResponse   *GetOIDCTokenResponse
	}{
		"password grant - success": {
			HttpResponse: GetOIDCTokenResponse{
				AccessToken:  "test-access-token",
				RefreshToken: "test-refresh-token",
				ExpiresIn:    300,
			},
			HttpResponseStatus: http.StatusOK,
			ExpectedError:      nil,
			Request: GetOIDCTokenRequest{
				Username:  "test-user",
				Password:  "test-pass",
				GrantType: "password",
			},
			ExpectedResponse: &GetOIDCTokenResponse{
				AccessToken:  "test-access-token",
				RefreshToken: "test-refresh-token",
				ExpiresIn:    300,
			},
		},
		"fail - invalid response": {
			HttpResponse:       []byte("invalid json"),
			HttpResponseStatus: http.StatusOK,
			Request: GetOIDCTokenRequest{
				Username:  "test-user",
				Password:  "test-pass",
				GrantType: "password",
			},
			ExpectedError:    fmt.Errorf("json: cannot unmarshal string into Go value of type keycloak.GetOIDCTokenResponse"),
			ExpectedResponse: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)
			cfg := &config.Config{
				KeycloakBaseURL:      "http://localhost:8080",
				KeycloakRealm:        "test-realm",
				KeycloakClientID:     "test-client",
				KeycloakClientSecret: "test-secret",
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/realms/test-realm/protocol/openid-connect/token", r.URL.Path)
				assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

				w.WriteHeader(tc.HttpResponseStatus)
				json.NewEncoder(w).Encode(tc.HttpResponse)
			}))
			defer ts.Close()
			cfg.KeycloakBaseURL = ts.URL

			client := &keycloakClient{
				httpClient: ts.Client(),
				config:     cfg,
				logger:     logger,
			}

			resp, err := client.GetOIDCToken(context.Background(), tc.Request)
			if tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.ExpectedResponse, resp)

		})
	}
}

func TestIntrospectOIDCToken(t *testing.T) {
	tests := map[string]struct {
		HttpResponse       any
		HttpResponseStatus int
		ExpectedError      error
		Request            IntrospectOIDCTokenRequest
		ExpectedResponse   *IntrospectOIDCTokenResponse
	}{
		"success": {
			HttpResponse: IntrospectOIDCTokenResponse{
				Active: true,
			},
			HttpResponseStatus: http.StatusOK,
			ExpectedError:      nil,
			Request: IntrospectOIDCTokenRequest{
				AccessToken: "access-token",
			},
			ExpectedResponse: &IntrospectOIDCTokenResponse{
				Active: true,
			},
		},
		"fail - invalid response": {
			HttpResponse:       []byte("invalid json"),
			HttpResponseStatus: http.StatusOK,
			Request: IntrospectOIDCTokenRequest{
				AccessToken: "access-token",
			},
			ExpectedError:    fmt.Errorf("json: cannot unmarshal string into Go value of type keycloak.IntrospectOIDCTokenResponse"),
			ExpectedResponse: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)
			cfg := &config.Config{
				KeycloakBaseURL:      "http://localhost:8080",
				KeycloakRealm:        "test-realm",
				KeycloakClientID:     "test-client",
				KeycloakClientSecret: "test-secret",
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/realms/test-realm/protocol/openid-connect/token/introspect", r.URL.Path)
				assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

				// Verify Basic Auth
				username, password, ok := r.BasicAuth()
				assert.True(t, ok)
				assert.Equal(t, "test-client", username)
				assert.Equal(t, "test-secret", password)

				w.WriteHeader(tc.HttpResponseStatus)
				json.NewEncoder(w).Encode(tc.HttpResponse)
			}))
			defer ts.Close()
			cfg.KeycloakBaseURL = ts.URL

			client := &keycloakClient{
				httpClient: ts.Client(),
				config:     cfg,
				logger:     logger,
			}

			resp, err := client.IntrospectOIDCToken(context.Background(), tc.Request)
			if tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.ExpectedResponse, resp)

		})
	}
}
