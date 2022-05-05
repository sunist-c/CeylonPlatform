package authentication

type TokenOptions struct {
	UserID   string
	ClientID string
	Token    string
	Scope    ScopeType
}

type ClientOptions struct {
	Name        string
	RedirectURL string
	Scope       ScopeType
	Method      AuthType
}

type UserOptions struct {
	Name     string
	Password string
	Scope    ScopeType
}

type AccessTokenOptions struct {
	UserID   string
	ClientID string
	Scope    ScopeType
}

type RefreshTokenOptions struct {
	UserID   string
	ClientID string
	Scope    ScopeType
}

type AuthorizationCodeOptions struct {
	UserID      string
	ClientID    string
	Scope       ScopeType
	RedirectURL string
}
