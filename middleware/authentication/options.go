package authentication

type TokenOptions struct {
}

type ClientOptions struct {
	id          string
	Name        string
	RedirectUri string
	Scope       ScopeType
	Method      AuthType
}

type UserOptions struct {
	id       string
	Name     string
	Password string
	Scope    ScopeType
}

type AccessTokenOptions struct {
	id       string
	UserID   string
	ClientID string
	Scope    ScopeType
}

type RefreshTokenOptions struct {
	id       string
	UserID   string
	ClientID string
	Scope    ScopeType
}

type AuthorizationCodeOptions struct {
	id          string
	UserID      string
	ClientID    string
	Scope       ScopeType
	RedirectUri string
}
