package authentication

import (
	"CeylonPlatform/middleware/initialization"
	"CeylonPlatform/pkg/uid"
	"errors"
	"github.com/go-redis/redis"
	"math/rand"
	"strconv"
	"time"
	"xorm.io/xorm"
)

type storageType string

const (
	Redis  storageType = "redis"
	Mysql  storageType = "mysql"
	Memory storageType = "memory"
	File   storageType = "file"
)

type Authenticator struct {
	dbConn           *xorm.Engine
	redisConn        *redis.Client
	tokenStorageType storageType
}

func DefaultAuthenticator() *Authenticator {
	return &Authenticator{
		dbConn:           initialization.DbConnection,
		redisConn:        initialization.RedisConnection,
		tokenStorageType: Redis,
	}
}

// SetStorageType 设置Token/RefreshToken/Code的存储方式，目前只有Redis实现
func (a *Authenticator) SetStorageType(storage storageType) {
	a.tokenStorageType = storage
}

func (a Authenticator) CreateClient() (client *Client, err error) {
	client = &Client{
		ID:          uid.GenerateUid(uid.Client),
		Name:        uid.GenerateMd5Len16(time.Now().String(), strconv.Itoa(rand.Int()%114514)),
		Key:         uid.GenerateUid(uid.ClientKey),
		Secret:      uid.GenerateUid(uid.ClientSecret),
		RedirectUri: "",
		Scope:       Student,
		Method:      PasswordAuth,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}

	_, err = a.dbConn.Insert(*client)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (a Authenticator) CreateClientWith(opts *ClientOptions) (client *Client, err error) {
	client = &Client{
		ID:          uid.GenerateUid(uid.Client),
		Name:        uid.GenerateMd5Len16(time.Now().String(), strconv.Itoa(rand.Int()%114514)),
		Key:         uid.GenerateUid(uid.ClientKey),
		Secret:      uid.GenerateUid(uid.ClientSecret),
		RedirectUri: "",
		Scope:       Student,
		Method:      PasswordAuth,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}

	if opts.Method != "" {
		client.Method = opts.Method
	}
	if opts.Name != "" {
		client.Name = opts.Name
	}
	if opts.RedirectUri != "" {
		client.RedirectUri = opts.RedirectUri
	}
	if opts.Scope != "" {
		client.Scope = opts.Scope
	}

	_, err = a.dbConn.Insert(*client)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (a Authenticator) UpdateClientWith(opts *ClientOptions) (client *Client, err error) {
	if opts.id == "" {
		return nil, errors.New("empty pk")
	}
}

func (a Authenticator) DeleteClient(opts *ClientOptions) (ok bool, err error) {

}

func (a Authenticator) CreateUser() (user *User, err error) {

}

func (a Authenticator) CreateUserWith(opts *UserOptions) (user *User, err error) {
}

func (a Authenticator) UpdateUserWith(opts *UserOptions) (user *User, err error) {
}

func (a Authenticator) DeleteUser(opts *UserOptions) (ok bool, err error) {

}

func (a Authenticator) CreateAccessToken() (token *AccessToken, err error) {

}

func (a Authenticator) CreateAccessTokenWith(opts *AccessTokenOptions) (token *AccessToken, err error) {

}

func (a Authenticator) UpdateAccessTokenWith(opts *AccessTokenOptions) (token *AccessToken, err error) {

}

func (a Authenticator) DeleteAccessToken(opts *AccessTokenOptions) (ok bool, err error) {

}

func (a Authenticator) CreateRefreshToken() (token *RefreshToken, err error) {

}

func (a Authenticator) CreateRefreshTokenWith(opts *RefreshTokenOptions) (token *RefreshToken, err error) {

}

func (a Authenticator) DeleteRefreshTokenWith(opts *RefreshTokenOptions) (ok bool, err error) {

}

func (a Authenticator) CreateAuthorizationCode() (token *AuthorizationCode, err error) {

}

func (a Authenticator) CreateAuthorizationCodeWith(opts *AuthorizationCodeOptions) (token *AuthorizationCode, err error) {

}

func (a Authenticator) UpdateAuthorizationCodeWith(opts *AuthorizationCodeOptions) (token *AuthorizationCode, err error) {

}

func (a Authenticator) DeleteAuthorizationCodeWith(opts *AuthorizationCodeOptions) (ok bool, err error) {

}

func (a Authenticator) PasswordAuth(userID, clientID, password string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error) {

}

func (a Authenticator) ClientAuth(clientID, clientSecret string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error) {

}

func (a Authenticator) ImplicitAuth(clientID, redirectUri string, scope ScopeType) (uri string, err error) {

}

func (a Authenticator) CodeAuth(clientID, redirectUri string, scope ScopeType) (uri string, err error) {

}

func (a Authenticator) CodeToToken(clientID, clientSecret, code string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error) {

}

func (a Authenticator) AuthToken(authType AuthType, opts *TokenOptions) (ok bool, err error) {

}

func (a Authenticator) RefreshToken(authType AuthType, opts *TokenOptions) (token *AccessToken, refreshToken *RefreshToken, err error) {

}
