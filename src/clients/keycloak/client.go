package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"go-api/src/config"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	_GET_OIDC_TOKEN_PATH        = "realms/%s/protocol/openid-connect/token"
	_INTROSPECT_OIDC_TOKEN_PATH = "realms/%s/protocol/openid-connect/token/introspect"
	_REVOKE_OIDC_TOKEN_PATH     = "realms/%s/protocol/openid-connect/revoke"

	_FORM_ENCODED = "application/x-www-form-urlencoded"
)

type KeycloakClient interface {
	GetOIDCToken(ctx context.Context, request GetOIDCTokenRequest) (*GetOIDCTokenResponse, error)
	RevokeOIDCToken(ctx context.Context, request RevokeOIDCTokenRequest) error
	IntrospectOIDCToken(ctx context.Context, request IntrospectOIDCTokenRequest) (*IntrospectOIDCTokenResponse, error)
}

type keycloakClient struct {
	httpClient *http.Client
	logger     *zap.Logger
	config     *config.Config
}

type KeycloakClientParams struct {
	fx.In

	Config *config.Config
	Logger *zap.Logger
}

func NewKeycloakClient(params KeycloakClientParams) KeycloakClient {
	httpClient := &http.Client{
		Timeout: time.Duration(params.Config.KeycloakTimoutMS) * time.Millisecond,
	}

	return &keycloakClient{
		httpClient: httpClient,
		config:     params.Config,
		logger:     params.Logger,
	}
}

func (c *keycloakClient) GetOIDCToken(ctx context.Context, request GetOIDCTokenRequest) (*GetOIDCTokenResponse, error) {
	if request.ClientID == "" {
		request.ClientID = c.config.KeycloakClientID
	}
	if request.ClientSecret == "" {
		request.ClientSecret = c.config.KeycloakClientSecret
	}

	path := fmt.Sprintf(_GET_OIDC_TOKEN_PATH, c.config.KeycloakRealm)
	body, err := buildFormEncodedBody(request)
	if err != nil {
		c.logger.Debug("failed to build form encoded body", zap.Error(err))
		return nil, err
	}

	params := requestParams{
		HTTPClient:  c.httpClient,
		BaseURL:     c.config.KeycloakBaseURL,
		Path:        path,
		ContentType: _FORM_ENCODED,
		Method:      http.MethodPost,
		Body:        body,
	}
	response, err := makeRequest[GetOIDCTokenResponse](params)
	if err != nil {
		c.logger.Debug("failed to make request", zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (c *keycloakClient) RevokeOIDCToken(ctx context.Context, request RevokeOIDCTokenRequest) error {
	if request.ClientID == "" {
		request.ClientID = c.config.KeycloakClientID
	}
	if request.ClientSecret == "" {
		request.ClientSecret = c.config.KeycloakClientSecret
	}

	path := fmt.Sprintf(_REVOKE_OIDC_TOKEN_PATH, c.config.KeycloakRealm)
	body, err := buildFormEncodedBody(request)
	if err != nil {
		c.logger.Debug("failed to build form encoded body", zap.Error(err))
		return err
	}
	params := requestParams{
		HTTPClient:  c.httpClient,
		BaseURL:     c.config.KeycloakBaseURL,
		Path:        path,
		ContentType: _FORM_ENCODED,
		Method:      http.MethodPost,
		Body:        body,
	}

	res, err := makeRequest[RevokeOIDCTokenResponse](params)
	c.logger.Debug("revoke oidc token response", zap.Any("response", res))

	if err != nil {
		c.logger.Debug("failed to make request", zap.Error(err))
		return err
	}
	return nil
}

func (c *keycloakClient) IntrospectOIDCToken(ctx context.Context, request IntrospectOIDCTokenRequest) (*IntrospectOIDCTokenResponse, error) {
	path := fmt.Sprintf(_INTROSPECT_OIDC_TOKEN_PATH, c.config.KeycloakRealm)
	body, err := buildFormEncodedBody(request)
	if err != nil {
		c.logger.Debug("failed to build form encoded body", zap.Error(err))
		return nil, err
	}

	params := requestParams{
		HTTPClient:  c.httpClient,
		BaseURL:     c.config.KeycloakBaseURL,
		Path:        path,
		ContentType: _FORM_ENCODED,
		Method:      http.MethodPost,
		Body:        body,
		Username:    c.config.KeycloakClientID,
		Password:    c.config.KeycloakClientSecret,
	}
	response, err := makeRequest[IntrospectOIDCTokenResponse](params)
	if err != nil {
		c.logger.Debug("failed to make request", zap.Error(err))
		return nil, err
	}
	if response == nil || !response.Active {
		return nil, fmt.Errorf("expired token")
	}
	return response, nil
}

func buildFormEncodedBody[T any](request T) (io.Reader, error) {
	formEnc, err := query.Values(&request)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(formEnc.Encode()), nil
}

type requestParams struct {
	HTTPClient  *http.Client
	BaseURL     string
	Path        string
	ContentType string
	Method      string
	Body        io.Reader
	Username    string
	Password    string
}

func makeRequest[T any](p requestParams) (*T, error) {
	baseUrl, err := url.Parse(p.BaseURL)
	if err != nil {
		return nil, err
	}
	url, err := baseUrl.Parse(p.Path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url.String(), p.Body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", p.ContentType)

	if p.Username != "" && p.Password != "" {
		req.SetBasicAuth(p.Username, p.Password)
	}

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(errorBody))
	}

	if resp.ContentLength == 0 {
		return nil, nil
	}

	var response T
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
