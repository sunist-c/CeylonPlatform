package authentication

import (
	"CeylonPlatform/middleware/initialization"
	"CeylonPlatform/pkg/uid"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
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

func generateRedisKey(uidType uid.UidType, id string) string {
	return fmt.Sprintf("%v:%v", uidType, id)
}

func encodePassword(password, id string) string {
	return uid.GenerateMd5Len32(password, id)
}

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
		Name:        uid.GenerateRandomName(),
		Key:         uid.GenerateUid(uid.ClientKey),
		Secret:      uid.GenerateUid(uid.ClientSecret),
		RedirectUri: "",
		Scope:       Student,
		Method:      PasswordAuth,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}

	_, err = a.dbConn.InsertOne(*client)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (a Authenticator) CreateClientWith(opts *ClientOptions) (client *Client, err error) {
	client = &Client{
		ID:          uid.GenerateUid(uid.Client),
		Name:        uid.GenerateRandomName(),
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
	if opts.Scope != Empty {
		client.Scope = opts.Scope
	}

	_, err = a.dbConn.InsertOne(*client)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (a Authenticator) UpdateClientWith(clientID string, opts *ClientOptions) (client *Client, err error) {
	client = &Client{}
	ok, err := a.dbConn.ID(clientID).Get(client)
	if !ok || err != nil {
		return nil, err
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
	if opts.Scope != Empty {
		client.Scope = opts.Scope
	}

	_, err = a.dbConn.ID(clientID).Update(*client)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (a Authenticator) DeleteClient(clientID string, opts *ClientOptions) (ok bool, err error) {
	client := &Client{
		ID: clientID,
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
	if opts.Scope != Empty {
		client.Scope = opts.Scope
	}

	_, err = a.dbConn.ID(clientID).Delete(*client)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (a Authenticator) CreateUser() (user *User, err error) {
	user = &User{
		ID:       uid.GenerateUid(uid.User),
		Name:     uid.GenerateRandomName(),
		Password: user.ID,
		Scope:    Student,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	_, err = a.dbConn.InsertOne(*user)
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (a Authenticator) CreateUserWith(opts *UserOptions) (user *User, err error) {
	user = &User{
		ID:       uid.GenerateUid(uid.User),
		Name:     uid.GenerateRandomName(),
		Password: user.Name,
		Scope:    Student,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	if opts.Name != "" {
		user.Name = opts.Name
	}
	if opts.Scope != Empty {
		user.Scope = opts.Scope
	}
	if opts.Password != "" {
		user.Password = uid.GenerateMd5Len32(opts.Password, user.ID)
	}

	_, err = a.dbConn.InsertOne(*user)
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (a Authenticator) UpdateUserWith(userID string, opts *UserOptions) (user *User, err error) {
	user = &User{}
	ok, err := a.dbConn.ID(userID).Get(user)
	if !ok || err != nil {
		return nil, err
	}

	if opts.Name != "" {
		user.Name = opts.Name
	}
	if opts.Password != "" {
		user.Password = encodePassword(opts.Password, userID)
	}
	if opts.Scope != Empty {
		user.Scope = opts.Scope
	}

	_, err = a.dbConn.ID(userID).Update(*user)
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (a Authenticator) DeleteUser(userID string, opts *UserOptions) (ok bool, err error) {
	user := &User{
		ID: userID,
	}
	if opts.Name != "" {
		user.Name = opts.Name
	}
	if opts.Password != "" {
		user.Password = encodePassword(opts.Password, userID)
	}
	if opts.Scope != Empty {
		user.Scope = opts.Scope
	}

	_, err = a.dbConn.ID(userID).Delete(*user)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (a Authenticator) CreateAccessTokenWith(opts *AccessTokenOptions) (token *AccessToken, err error) {
	if opts.UserID == "" || opts.ClientID == "" || opts.Scope == Empty {
		return nil, errors.New("bad options with empty field")
	}

	token = &AccessToken{
		ID:       uid.GenerateUid(uid.Token),
		Token:    uid.GenerateUid(uid.Token),
		UserID:   opts.UserID,
		ClientID: opts.ClientID,
		Scope:    opts.Scope,
		ExpireAt: time.Now().Add(time.Second * 3600),
		CreateAt: time.Now(),
	}

	err = a.redisConn.Set(generateRedisKey(uid.Token, token.ID), *token, time.Second*3600).Err()
	if err != nil {
		return nil, err
	} else {
		return token, nil
	}
}

func (a Authenticator) CreateRefreshTokenWith(opts *RefreshTokenOptions) (token *RefreshToken, err error) {
	if opts.UserID == "" || opts.ClientID == "" || opts.Scope == Empty {
		return nil, errors.New("bad options with empty field")
	}

	token = &RefreshToken{
		ID:       uid.GenerateUid(uid.RefreshToken),
		Token:    uid.GenerateUid(uid.RefreshToken),
		UserID:   opts.UserID,
		ClientID: opts.ClientID,
		Scope:    opts.Scope,
		ExpireAt: time.Now().Add(time.Second * 3600),
		CreateAt: time.Now(),
	}

	err = a.redisConn.Set(generateRedisKey(uid.RefreshToken, token.ID), *token, time.Second*3600).Err()
	if err != nil {
		return nil, err
	} else {
		return token, nil
	}
}

func (a Authenticator) CreateAuthorizationCodeWith(opts *AuthorizationCodeOptions) (code *AuthorizationCode, err error) {
	if opts.UserID == "" || opts.ClientID == "" || opts.Scope == Empty {
		return nil, errors.New("bad options with empty field")
	}

	code = &AuthorizationCode{
		ID:       uid.GenerateUid(uid.AuthCode),
		Code:     uid.GenerateUid(uid.AuthCode),
		UserID:   opts.UserID,
		ClientID: opts.ClientID,
		Scope:    opts.Scope,
		ExpireAt: time.Now().Add(time.Second * 3600),
		CreateAt: time.Now(),
	}

	err = a.redisConn.Set(generateRedisKey(uid.AuthCode, code.ID), *code, time.Second*3600).Err()
	if err != nil {
		return nil, err
	} else {
		return code, nil
	}
}

func (a Authenticator) PasswordAuth(userID, clientID, password string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error) {
	client := &Client{}
	user := &User{}

	ok, err := a.dbConn.ID(clientID).Get(&client)
	if !ok || err != nil {
		return nil, nil, err
	}
	ok, err = a.dbConn.ID(userID).Get(&user)
	if !ok || err != nil {
		return nil, nil, err
	}

	if user.Scope < scope || client.Scope < scope {
		return nil, nil, errors.New("bad auth with illegal scope")
	}
	if user.Password != encodePassword(password, userID) {
		return nil, nil, errors.New("bad auth with incorrect password")
	}

	token, err = a.CreateAccessTokenWith(&AccessTokenOptions{
		UserID:   userID,
		ClientID: clientID,
		Scope:    scope,
	})
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err = a.CreateRefreshTokenWith(&RefreshTokenOptions{
		UserID:   userID,
		ClientID: clientID,
		Scope:    scope,
	})
	if err != nil {
		return nil, nil, err
	}

	return token, refreshToken, nil
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
