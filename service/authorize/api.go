package authorize

import (
	"CeylonPlatform/middleware/api"
	"CeylonPlatform/middleware/authentication"
	"github.com/gin-gonic/gin"
	"time"
)

func init() {
	api.AddService(Service{})
}

type LoginPageViewModel struct {
	User struct {
		Name     string
		Password string
	}
	Client struct {
		Name    string
		ID      string
		Website string
		Favicon string
		Scope   authentication.ScopeType
	}
}

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

type Service struct {
}

func (s Service) Handlers() []*api.ServiceHandler {
	generateToken := api.ServiceHandler{
		Name:   "oauth-api",
		Url:    "oauth",
		Method: api.POST,
		Handler: func(req *gin.Context, ctx *api.Context) {
			switch authentication.AuthType(req.Query("auth_type")) {
			case authentication.PasswordAuth:
				request := &passwordAuthRequest{}
				err := req.ShouldBindJSON(request)
				if err != nil {
					api.BadRequestError(req)
					return
				}

				token, refreshToken, err := authentication.DefaultAuthenticator().PasswordAuth(request.UserID, request.ClientID, request.Password, request.Scope)
				if err != nil || token == nil || refreshToken == nil {
					api.InternalServerError(req)
					return
				}

				req.JSON(200, passwordAuthResponse{
					Token:        token.Token,
					RefreshToken: refreshToken.Token,
					ExpireAt:     token.ExpireAt,
				})
			case authentication.ClientAuth:
				request := &clientAuthRequest{}
				err := req.ShouldBindJSON(request)
				if err != nil {
					api.BadRequestError(req)
					return
				}

				token, refreshToken, err := authentication.DefaultAuthenticator().ClientAuth(request.ClientID, request.ClientSecret, request.Scope)
				if err != nil || token == nil || refreshToken == nil {
					api.InternalServerError(req)
					return
				}

				req.JSON(200, clientAuthResponse{
					Token:        token.Token,
					RefreshToken: refreshToken.Token,
					ExpireAt:     token.ExpireAt,
				})
			case authentication.ImplicitAuth:
			case authentication.CodeAuth:
			default:
				api.BadRequestError(req)
				return
			}
		},
		Dependence: nil,
		Router:     api.BaseRouter().Group("/authorize"),
	}

	authToken := api.ServiceHandler{
		Name:   "auth-api",
		Url:    "introspect",
		Method: api.GET,
		Handler: func(req *gin.Context, ctx *api.Context) {
			switch authentication.AuthType(req.Query("auth_type")) {
			case authentication.ClientAuth:
				request := &clientTokenAuthRequest{}
				err := req.ShouldBindJSON(req)
				if err != nil {
					api.BadRequestError(req)
					return
				}

				ok, err := authentication.DefaultAuthenticator().AuthToken(authentication.ClientAuth, &authentication.TokenOptions{
					ClientID: request.ClientID,
					Token:    request.Token,
					Scope:    request.Scope,
				})

				req.JSON(200, clientTokenAuthResponse{Status: ok})
			case authentication.PasswordAuth:
				request := &passwordTokenAuthRequest{}
				err := req.ShouldBindJSON(request)
				if err != nil {
					api.BadRequestError(req)
					return
				}

				ok, err := authentication.DefaultAuthenticator().AuthToken(authentication.PasswordAuth, &authentication.TokenOptions{
					ClientID: request.ClientID,
					UserID:   request.UserID,
					Token:    request.Token,
					Scope:    request.Scope,
				})

				req.JSON(200, passwordTokenAuthResponse{Status: ok})
			case authentication.CodeAuth:

			case authentication.ImplicitAuth:

			default:
				api.BadRequestError(req)
				return
			}
		},
		Dependence: nil,
		Router:     api.BaseRouter().Group("/authorize"),
	}

	return []*api.ServiceHandler{
		&generateToken,
		&authToken,
	}
}
