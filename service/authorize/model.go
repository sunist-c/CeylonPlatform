package authorize

import (
	"CeylonPlatform/middleware/authentication"
	"time"
)

type passwordAuthRequest struct {
	ClientID string                   `json:"client_id"`
	UserID   string                   `json:"user_id"`
	Password string                   `json:"password"`
	Scope    authentication.ScopeType `json:"scope"`
}

type passwordAuthResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpireAt     time.Time `json:"expire_at"`
}

type clientAuthRequest struct {
	ClientID     string                   `json:"client_id"`
	ClientSecret string                   `json:"client_secret"`
	Scope        authentication.ScopeType `json:"scope"`
}

type clientAuthResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpireAt     time.Time `json:"expire_at"`
}

type passwordTokenAuthRequest struct {
	Token    string                   `json:"token"`
	ClientID string                   `json:"client_id"`
	UserID   string                   `json:"user_id"`
	Scope    authentication.ScopeType `json:"scope"`
}

type passwordTokenAuthResponse struct {
	Status bool `json:"status"`
}

type clientTokenAuthRequest struct {
	Token    string                   `json:"token"`
	ClientID string                   `json:"client_id"`
	Scope    authentication.ScopeType `json:"scope"`
}

type clientTokenAuthResponse struct {
	Status bool `json:"status"`
}

type registerRequest struct {
	UserList []struct {
		Username string                   `json:"username"`
		Password string                   `json:"password"`
		Scope    authentication.ScopeType `json:"scope"`
	} `json:"user_list"`
}

type registerResponse struct {
	UserList []struct {
		UserID   string                   `json:"user_id"`
		Username string                   `json:"username"`
		Scope    authentication.ScopeType `json:"scope"`
		CreateAt time.Time                `json:"create_at"`
	} `json:"user_list"`
}

type registerClientRequest struct {
	ClientName   string                   `json:"client_name"`
	ClientDomain string                   `json:"client_domain"`
	Scope        authentication.ScopeType `json:"scope"`
	AuthType     authentication.AuthType  `json:"auth_type"`
}

type registerClientResponse struct {
	ClientID     string                   `json:"client_id"`
	ClientKey    string                   `json:"client_key"`
	ClientSecret string                   `json:"client_secret"`
	ClientDomain string                   `json:"client_domain"`
	Scope        authentication.ScopeType `json:"scope"`
	AuthType     authentication.AuthType  `json:"auth_type"`
	CreateAt     time.Time                `json:"create_at"`
}

type loginPageViewModel struct {
	Favicon       string
	ClientName    string
	ClientWebsite string
	RedirectUrl   string
	AuthType      authentication.AuthType
	Scope         authentication.ScopeType
	Username      string
	Password      string
	ClientID      string
	Status        string
}
