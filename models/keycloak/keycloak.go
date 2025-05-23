package keycloak

type GetOIDCTokenRequest struct {
	GrantType    string `url:"grant_type"`
	Scope        string `url:"scope,omitempty"`
	Username     string `url:"username,omitempty"`
	Password     string `url:"password,omitempty"`
	ClientID     string `url:"client_id"`
	ClientSecret string `url:"client_secret,omitempty"`
	RefreshToken string `url:"refresh_token,omitempty"`
}

type RevokeOIDCTokenRequest struct {
	Token        string `url:"token"`
	ClientID     string `url:"client_id"`
	ClientSecret string `url:"client_secret,omitempty"`
}

type RevokeOIDCTokenResponse struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type IntrospectOIDCTokenRequest struct {
	AccessToken string `url:"token"`
}

type GetOIDCTokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not_before_policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	Account Account `json:"account"`
}

type Account struct {
	Roles []string `json:"roles"`
}

type IntrospectOIDCTokenResponse struct {
	Exp               int            `json:"exp"`
	Iat               int            `json:"iat"`
	Jti               string         `json:"jti"`
	Iss               string         `json:"iss"`
	Aud               string         `json:"aud"`
	Sub               string         `json:"sub"`
	Azp               string         `json:"azp"`
	Sid               string         `json:"sid"`
	Acr               string         `json:"acr"`
	AllowedOrigins    []string       `json:"allowed-origins"`
	RealmAccess       RealmAccess    `json:"realm_access"`
	ResourceAccess    ResourceAccess `json:"resource_access"`
	Scope             string         `json:"scope"`
	EmailVerified     bool           `json:"email_verified"`
	Name              string         `json:"name"`
	PreferredUsername string         `json:"preferred_username"`
	GivenName         string         `json:"given_name"`
	FamilyName        string         `json:"family_name"`
	Email             string         `json:"email"`
	ClientID          string         `json:"client_id"`
	Username          string         `json:"username"`
	TokenType         string         `json:"token_type"`
	Active            bool           `json:"active"`
}
