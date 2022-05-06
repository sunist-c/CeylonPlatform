package authorize

import (
	"CeylonPlatform/middleware/api"
	"CeylonPlatform/middleware/authentication"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func init() {
	api.AddService(Service{})
}

type Service struct {
}

func (s Service) Handlers() []*api.ServiceHandler {
	router := api.BaseRouter().Group("/authorize")
	auth := authentication.DefaultAuthenticator()

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

				token, refreshToken, err := auth.PasswordAuth(request.UserID, request.ClientID, request.Password, request.Scope)
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

				token, refreshToken, err := auth.ClientAuth(request.ClientID, request.ClientSecret, request.Scope)
				if err != nil || token == nil || refreshToken == nil {
					api.InternalServerError(req)
					return
				}

				req.JSON(200, clientAuthResponse{
					Token:        token.Token,
					RefreshToken: refreshToken.Token,
					ExpireAt:     token.ExpireAt,
				})
			case authentication.CodeAuth:
				username := req.PostForm("username")
				password := req.PostForm("password")
				redirectURL := req.Query("redirect_url")
				clientID := req.Query("client_id")
				scopeStr := req.Query("scope")
				scope, err := strconv.Atoi(scopeStr)
				if clientID == "" || err != nil || redirectURL == "" {
					api.BadRequestError(req)
					return
				}

				userID, err := auth.GetUserID(username)
				if err != nil {
					// todo: render html
					req.String(403, "用户名错误")
					return
				}

				url, err := auth.CodeAuth(userID, clientID, redirectURL, password, authentication.ScopeType(scope))
				if err != nil {
					// todo: render html
					req.String(403, "权限验证失败")
					return
				}

				req.Redirect(302, url)
			case authentication.ImplicitAuth:
				username := req.PostForm("username")
				password := req.PostForm("password")
				redirectURL := req.Query("redirect_url")
				clientID := req.Query("client_id")
				scopeStr := req.Query("scope")
				scope, err := strconv.Atoi(scopeStr)
				if clientID == "" || err != nil || redirectURL == "" {
					api.BadRequestError(req)
					return
				}

				userID, err := auth.GetUserID(username)
				if err != nil {
					// todo: render html
					req.String(403, "用户名错误")
					return
				}

				url, err := auth.ImplicitAuth(userID, clientID, redirectURL, password, authentication.ScopeType(scope))
				if err != nil {
					// todo: render html
					req.String(403, "权限验证失败")
					return
				}

				req.Redirect(302, url)
			default:
				api.BadRequestError(req)
				return
			}
		},
		Dependence: nil,
		Router:     router,
	}

	authToken := api.ServiceHandler{
		Name:   "auth-api",
		Url:    "introspect",
		Method: api.GET,
		Handler: func(req *gin.Context, ctx *api.Context) {
			switch authentication.AuthType(req.Query("auth_type")) {
			case authentication.ClientAuth:
				request := &clientTokenAuthRequest{}
				err := req.ShouldBindJSON(request)
				if err != nil {
					api.BadRequestError(req)
					return
				}

				ok, err := auth.AuthToken(authentication.ClientAuth, &authentication.TokenOptions{
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

				ok, err := auth.AuthToken(authentication.PasswordAuth, &authentication.TokenOptions{
					ClientID: request.ClientID,
					UserID:   request.UserID,
					Token:    request.Token,
					Scope:    request.Scope,
				})

				req.JSON(200, passwordTokenAuthResponse{Status: ok})
			case authentication.CodeAuth:
				request := &passwordTokenAuthRequest{}
				err := req.ShouldBindJSON(request)
				if err != nil {
					api.BadRequestError(req)
					return
				}

				ok, err := auth.AuthToken(authentication.PasswordAuth, &authentication.TokenOptions{
					ClientID: request.ClientID,
					UserID:   request.UserID,
					Token:    request.Token,
					Scope:    request.Scope,
				})

				req.JSON(200, passwordTokenAuthResponse{Status: ok})
			case authentication.ImplicitAuth:
				request := &passwordTokenAuthRequest{}
				err := req.ShouldBindJSON(request)
				if err != nil {
					api.BadRequestError(req)
					return
				}

				ok, err := auth.AuthToken(authentication.PasswordAuth, &authentication.TokenOptions{
					ClientID: request.ClientID,
					UserID:   request.UserID,
					Token:    request.Token,
					Scope:    request.Scope,
				})

				req.JSON(200, passwordTokenAuthResponse{Status: ok})
			default:
				api.BadRequestError(req)
				return
			}
		},
		Dependence: nil,
		Router:     router,
	}

	loginHandler := api.ServiceHandler{
		Name:   "login-api",
		Url:    "login",
		Method: api.GET,
		Handler: func(req *gin.Context, ctx *api.Context) {
			authType := req.Query("auth_type")
			clientID := req.Query("client_id")
			scopeStr := req.Query("scope")
			status := req.Query("status")
			scope, err := strconv.Atoi(scopeStr)
			redirectURL := req.Query("redirect_url")
			if authType == "" || clientID == "" || err != nil || redirectURL == "" {
				api.BadRequestError(req)
				return
			}

			client, err := auth.GetClientInfo(clientID)
			if err != nil {
				api.BadRequestError(req)
				return
			}

			req.HTML(200, "authorize/login.tmpl", loginPageViewModel{
				Favicon:       "https://sunist.cn/assets/img/favicon.png",
				ClientID:      clientID,
				ClientName:    client.Name,
				ClientWebsite: client.RedirectDomain,
				RedirectUrl:   redirectURL,
				AuthType:      authentication.AuthType(authType),
				Scope:         authentication.ScopeType(scope),
				Status:        status,
			})
		},
		Dependence: nil,
		Router:     router,
	}

	registerHandler := api.ServiceHandler{
		Name:   "register-api",
		Url:    "users",
		Method: api.POST,
		Handler: func(req *gin.Context, ctx *api.Context) {
			request := registerRequest{}
			err := req.ShouldBindJSON(&request)
			if err != nil {
				api.BadRequestError(req)
				return
			}

			response := registerResponse{
				UserList: make([]struct {
					UserID   string                   `json:"user_id"`
					Username string                   `json:"username"`
					Scope    authentication.ScopeType `json:"scope"`
					CreateAt time.Time                `json:"create_at"`
				}, 0, 1),
			}
			for i, _ := range request.UserList {
				user, err := auth.CreateUserWith(&authentication.UserOptions{
					Name:     request.UserList[i].Username,
					Password: request.UserList[i].Password,
					Scope:    request.UserList[i].Scope,
				})
				if err != nil {
					req.JSON(500, response)
					return
				} else {
					response.UserList = append(response.UserList, struct {
						UserID   string                   `json:"user_id"`
						Username string                   `json:"username"`
						Scope    authentication.ScopeType `json:"scope"`
						CreateAt time.Time                `json:"create_at"`
					}{UserID: user.ID, Username: user.Name, Scope: user.Scope, CreateAt: user.CreateAt})
				}
			}
			req.JSON(200, response)
		},
		Dependence: nil,
		Router:     router,
	}

	registerClientHandler := api.ServiceHandler{
		Name:   "register-client-handler",
		Url:    "client",
		Method: api.POST,
		Handler: func(req *gin.Context, ctx *api.Context) {
			request := registerClientRequest{}
			err := req.ShouldBindJSON(&request)
			if err != nil {
				api.BadRequestError(req)
				return
			}

			client, err := auth.CreateClientWith(&authentication.ClientOptions{
				Name:        request.ClientName,
				RedirectURL: request.ClientDomain,
				Scope:       request.Scope,
				Method:      request.AuthType,
			})
			if err != nil {
				api.InternalServerError(req)
				return
			} else {
				req.JSON(200, registerClientResponse{
					ClientID:     client.ID,
					ClientKey:    client.Key,
					ClientSecret: client.Secret,
					ClientDomain: client.RedirectDomain,
					Scope:        client.Scope,
					AuthType:     client.Method,
					CreateAt:     client.CreateAt,
				})
			}
		},
		Dependence: nil,
		Router:     router,
	}

	return []*api.ServiceHandler{
		&generateToken,
		&authToken,
		&loginHandler,
		&registerHandler,
		&registerClientHandler,
	}
}
