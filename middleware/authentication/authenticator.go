package authentication

import (
	"CeylonPlatform/middleware/initialization"
	"CeylonPlatform/pkg/uid"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
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
		ID:             uid.GenerateUid(uid.Client),
		Name:           uid.GenerateRandomName(),
		Key:            uid.GenerateUid(uid.ClientKey),
		Secret:         uid.GenerateUid(uid.ClientSecret),
		RedirectDomain: "",
		Scope:          Student,
		Method:         PasswordAuth,
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
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
		ID:             uid.GenerateUid(uid.Client),
		Name:           uid.GenerateRandomName(),
		Key:            uid.GenerateUid(uid.ClientKey),
		Secret:         uid.GenerateUid(uid.ClientSecret),
		RedirectDomain: "",
		Scope:          Student,
		Method:         PasswordAuth,
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
	}

	if opts.Method != "" {
		client.Method = opts.Method
	}
	if opts.Name != "" {
		client.Name = opts.Name
	}
	if opts.RedirectURL != "" {
		client.RedirectDomain = opts.RedirectURL
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
	if opts.RedirectURL != "" {
		client.RedirectDomain = opts.RedirectURL
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
	if opts.RedirectURL != "" {
		client.RedirectDomain = opts.RedirectURL
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
	userID := uid.GenerateUid(uid.User)
	user = &User{
		ID:       userID,
		Name:     uid.GenerateRandomName(),
		Password: userID,
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
	userID := uid.GenerateUid(uid.User)
	user = &User{
		ID:       userID,
		Name:     uid.GenerateRandomName(),
		Password: userID,
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
	if opts.ClientID == "" || opts.Scope == Empty {
		return nil, errors.New("bad options with empty field in necessary fields")
	}

	tokenStr := uid.GenerateUid(uid.Token)

	token = &AccessToken{
		ID:       tokenStr,
		Token:    tokenStr,
		UserID:   opts.UserID,
		ClientID: opts.ClientID,
		Scope:    opts.Scope,
		ExpireAt: time.Now().Add(time.Second * 3600),
		CreateAt: time.Now(),
	}

	redisObjStr, err := json.Marshal(*token)
	if err != nil {
		return nil, err
	}

	err = a.redisConn.Set(generateRedisKey(uid.Token, token.ID), string(redisObjStr), time.Second*3600).Err()
	if err != nil {
		return nil, err
	} else {
		return token, nil
	}
}

func (a Authenticator) CreateRefreshTokenWith(opts *RefreshTokenOptions) (token *RefreshToken, err error) {
	if opts.ClientID == "" || opts.Scope == Empty {
		return nil, errors.New("bad options with empty field in necessary fields")
	}

	tokenStr := uid.GenerateUid(uid.RefreshToken)

	token = &RefreshToken{
		ID:       tokenStr,
		Token:    tokenStr,
		UserID:   opts.UserID,
		ClientID: opts.ClientID,
		Scope:    opts.Scope,
		ExpireAt: time.Now().Add(time.Second * 3600),
		CreateAt: time.Now(),
	}

	redisObjStr, err := json.Marshal(*token)
	if err != nil {
		return nil, err
	}

	err = a.redisConn.Set(generateRedisKey(uid.RefreshToken, token.ID), string(redisObjStr), time.Second*3600).Err()
	if err != nil {
		return nil, err
	} else {
		return token, nil
	}
}

func (a Authenticator) CreateAuthorizationCodeWith(opts *AuthorizationCodeOptions) (code *AuthorizationCode, err error) {
	if opts.ClientID == "" || opts.Scope == Empty {
		return nil, errors.New("bad options with empty field in necessary fields")
	}

	codeStr := uid.GenerateUid(uid.AuthCode)

	code = &AuthorizationCode{
		ID:       codeStr,
		Code:     codeStr,
		UserID:   opts.UserID,
		ClientID: opts.ClientID,
		Scope:    opts.Scope,
		ExpireAt: time.Now().Add(time.Second * 3600),
		CreateAt: time.Now(),
	}

	redisObjStr, err := json.Marshal(*code)
	if err != nil {
		return nil, err
	}

	err = a.redisConn.Set(generateRedisKey(uid.AuthCode, code.ID), redisObjStr, time.Second*3600).Err()
	if err != nil {
		return nil, err
	} else {
		return code, nil
	}
}

func (a Authenticator) PasswordAuth(userID, clientID, password string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error) {
	client := &Client{}
	user := &User{}

	ok, err := a.dbConn.ID(clientID).Get(client)
	if !ok || err != nil {
		return nil, nil, err
	}
	ok, err = a.dbConn.ID(userID).Get(user)
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
	client := &Client{}
	ok, err := a.dbConn.ID(clientID).Get(client)
	if !ok || err != nil {
		return nil, nil, err
	}

	if client.Scope < scope {
		return nil, nil, errors.New("bad auth with illegal scope")
	}

	if client.Secret != clientSecret {
		return nil, nil, errors.New("bad auth with incorrect secret")
	}

	token, err = a.CreateAccessTokenWith(&AccessTokenOptions{
		ClientID: clientID,
		Scope:    scope,
	})
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err = a.CreateRefreshTokenWith(&RefreshTokenOptions{
		ClientID: clientID,
		Scope:    scope,
	})
	if err != nil {
		return nil, nil, err
	}

	return token, refreshToken, nil
}

// ImplicitAuth Client发起Implicit认证后，用户在页面授权完成后，重定向到Client的业务逻辑
func (a Authenticator) ImplicitAuth(userID, clientID, redirectURL string, scope ScopeType) (uri string, err error) {
	client := &Client{}
	user := &User{}
	ok, err := a.dbConn.ID(clientID).Get(client)
	if !ok || err != nil {
		return "", err
	}
	ok, err = a.dbConn.ID(userID).Get(user)
	if !ok || err != nil {
		return "", err
	}

	if !strings.Contains(redirectURL, client.RedirectDomain) {
		return "", errors.New("bad auth with illegal redirect-url")
	}

	if client.Scope < scope || user.Scope < scope {
		return "", errors.New("bad auth with illegal scope")
	}

	token, err := a.CreateAccessTokenWith(&AccessTokenOptions{
		UserID:   userID,
		ClientID: clientID,
		Scope:    scope,
	})
	if err != nil {
		return "", err
	}

	refreshToken, err := a.CreateRefreshTokenWith(&RefreshTokenOptions{
		UserID:   userID,
		ClientID: clientID,
		Scope:    scope,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v?access-token=%v&refresh-token=%v", redirectURL, token.Token, refreshToken.Token), nil
}

// CodeAuth Client发起Code认证后，用户在页面授权完成后，重定向到Client后端，供Client后端进行CodeToToken的业务逻辑
func (a Authenticator) CodeAuth(userID, clientID, redirectURL string, scope ScopeType) (uri string, err error) {
	client := &Client{}
	user := &User{}
	ok, err := a.dbConn.ID(clientID).Get(client)
	if !ok || err != nil {
		return "", err
	}
	ok, err = a.dbConn.ID(userID).Get(user)
	if !ok || err != nil {
		return "", err
	}

	if !strings.Contains(redirectURL, client.RedirectDomain) {
		return "", errors.New("bad auth with illegal redirect-url")
	}

	if client.Scope < scope || user.Scope < scope {
		return "", errors.New("bad auth with illegal scope")
	}

	code, err := a.CreateAuthorizationCodeWith(&AuthorizationCodeOptions{
		UserID:      userID,
		ClientID:    clientID,
		Scope:       scope,
		RedirectURL: redirectURL,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v?code=%v", redirectURL, code.Code), nil
}

// CodeToToken Code认证模式下，Client的后端处理了Code回调，并请求CodeToToken步骤的业务逻辑
func (a Authenticator) CodeToToken(userID, clientID, clientSecret, code string, scope ScopeType) (token *AccessToken, refreshToken *RefreshToken, err error) {
	client := &Client{}
	user := &User{}
	authCode := &AuthorizationCode{}
	ok, err := a.dbConn.ID(clientID).Get(client)
	if !ok || err != nil {
		return nil, nil, err
	}
	ok, err = a.dbConn.ID(userID).Get(user)
	if !ok || err != nil {
		return nil, nil, err
	}

	if client.Scope < scope || user.Scope < scope {
		return nil, nil, errors.New("bad auth with illegal scope")
	}

	if client.Secret != clientSecret {
		return nil, nil, errors.New("bad auth with incorrect secret")
	}

	authCodeStr, err := a.redisConn.Get(generateRedisKey(uid.AuthCode, code)).Result()
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal([]byte(authCodeStr), authCode)
	if err != nil {
		return nil, nil, err
	}

	if authCode.Code != code {
		return nil, nil, errors.New("bad auth with incorrect code")
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

// AuthToken 验证Token是否正确
func (a Authenticator) AuthToken(authType AuthType, opts *TokenOptions) (ok bool, err error) {
	switch authType {
	case ClientAuth:
		if opts.Token == "" || opts.ClientID == "" || opts.Scope == Empty {
			return false, errors.New("bad options with empty field in necessary fields")
		}

		tokenStr, err := a.redisConn.Get(generateRedisKey(uid.Token, opts.Token)).Result()
		if err != nil {
			return false, err
		}
		token := &AccessToken{}
		err = json.Unmarshal([]byte(tokenStr), token)
		if err != nil {
			return false, err
		}

		if token.Scope < opts.Scope || token.ClientID != opts.ClientID || token.Token != opts.Token {
			return false, nil
		} else {
			return true, nil
		}
	case PasswordAuth:
		if opts.Token == "" || opts.UserID == "" || opts.ClientID == "" || opts.Scope == Empty {
			return false, errors.New("bad options with empty field in necessary fields")
		}

		tokenStr, err := a.redisConn.Get(generateRedisKey(uid.Token, opts.Token)).Result()
		if err != nil {
			return false, err
		}
		token := &AccessToken{}
		err = json.Unmarshal([]byte(tokenStr), token)
		if err != nil {
			return false, err
		}

		if token.Scope < opts.Scope || token.ClientID != opts.ClientID || token.Token != opts.Token || token.UserID != opts.UserID {
			return false, nil
		} else {
			return true, nil
		}
	case CodeAuth:
		if opts.Token == "" || opts.UserID == "" || opts.ClientID == "" || opts.Scope == Empty {
			return false, errors.New("bad options with empty field in necessary fields")
		}

		tokenStr, err := a.redisConn.Get(generateRedisKey(uid.Token, opts.Token)).Result()
		if err != nil {
			return false, err
		}
		token := &AccessToken{}
		err = json.Unmarshal([]byte(tokenStr), token)
		if err != nil {
			return false, err
		}

		if token.Scope < opts.Scope || token.ClientID != opts.ClientID || token.Token != opts.Token || token.UserID != opts.UserID {
			return false, nil
		} else {
			return true, nil
		}
	case ImplicitAuth:
		if opts.Token == "" || opts.UserID == "" || opts.ClientID == "" || opts.Scope == Empty {
			return false, errors.New("bad options with empty field in necessary fields")
		}

		tokenStr, err := a.redisConn.Get(generateRedisKey(uid.Token, opts.Token)).Result()
		if err != nil {
			return false, err
		}
		token := &AccessToken{}
		err = json.Unmarshal([]byte(tokenStr), token)
		if err != nil {
			return false, err
		}

		if token.Scope < opts.Scope || token.ClientID != opts.ClientID || token.Token != opts.Token || token.UserID != opts.UserID {
			return false, nil
		} else {
			return true, nil
		}
	default:
		return false, errors.New("bad auth-type")
	}
}

// AuthRefreshToken 验证RefreshToken是否正确
func (a Authenticator) AuthRefreshToken(authType AuthType, opts *TokenOptions) (ok bool, err error) {
	switch authType {
	case ClientAuth:
		if opts.Token == "" || opts.ClientID == "" || opts.Scope == Empty {
			return false, errors.New("bad options with empty field in necessary fields")
		}

		tokenStr, err := a.redisConn.Get(generateRedisKey(uid.RefreshToken, opts.Token)).Result()
		if err != nil {
			return false, err
		}
		token := &RefreshToken{}
		err = json.Unmarshal([]byte(tokenStr), token)
		if err != nil {
			return false, err
		}

		if token.Scope < opts.Scope || token.ClientID != opts.ClientID || token.Token != opts.Token {
			return false, nil
		} else {
			return true, nil
		}
	case PasswordAuth:
		if opts.Token == "" || opts.UserID == "" || opts.ClientID == "" || opts.Scope == Empty {
			return false, errors.New("bad options with empty field in necessary fields")
		}

		tokenStr, err := a.redisConn.Get(generateRedisKey(uid.RefreshToken, opts.Token)).Result()
		if err != nil {
			return false, err
		}
		token := &RefreshToken{}
		err = json.Unmarshal([]byte(tokenStr), token)
		if err != nil {
			return false, err
		}

		if token.Scope < opts.Scope || token.ClientID != opts.ClientID || token.Token != opts.Token || token.UserID != opts.UserID {
			return false, nil
		} else {
			return true, nil
		}
	case CodeAuth:
		if opts.Token == "" || opts.UserID == "" || opts.ClientID == "" || opts.Scope == Empty {
			return false, errors.New("bad options with empty field in necessary fields")
		}

		tokenStr, err := a.redisConn.Get(generateRedisKey(uid.RefreshToken, opts.Token)).Result()
		if err != nil {
			return false, err
		}
		token := &RefreshToken{}
		err = json.Unmarshal([]byte(tokenStr), token)
		if err != nil {
			return false, err
		}

		if token.Scope < opts.Scope || token.ClientID != opts.ClientID || token.Token != opts.Token || token.UserID != opts.UserID {
			return false, nil
		} else {
			return true, nil
		}
	case ImplicitAuth:
		if opts.Token == "" || opts.UserID == "" || opts.ClientID == "" || opts.Scope == Empty {
			return false, errors.New("bad options with empty field in necessary fields")
		}

		tokenStr, err := a.redisConn.Get(generateRedisKey(uid.RefreshToken, opts.Token)).Result()
		if err != nil {
			return false, err
		}
		token := &RefreshToken{}
		err = json.Unmarshal([]byte(tokenStr), token)
		if err != nil {
			return false, err
		}

		if token.Scope < opts.Scope || token.ClientID != opts.ClientID || token.Token != opts.Token || token.UserID != opts.UserID {
			return false, nil
		} else {
			return true, nil
		}
	default:
		return false, errors.New("bad auth-type")
	}
}

// RefreshToken 使用RefreshToken刷新Token与RefreshToken
func (a Authenticator) RefreshToken(authType AuthType, opts *TokenOptions) (token *AccessToken, refreshToken *RefreshToken, err error) {
	switch authType {
	case ClientAuth:
		ok, err := a.AuthRefreshToken(ClientAuth, opts)
		if !ok || err != nil {
			return nil, nil, err
		}

		token, err = a.CreateAccessTokenWith(&AccessTokenOptions{
			ClientID: opts.ClientID,
			Scope:    opts.Scope,
		})
		if err != nil {
			return nil, nil, err
		}

		refreshToken, err = a.CreateRefreshTokenWith(&RefreshTokenOptions{
			ClientID: opts.ClientID,
			Scope:    opts.Scope,
		})
		if err != nil {
			return nil, nil, err
		}

		return token, refreshToken, nil
	case CodeAuth:
		ok, err := a.AuthRefreshToken(CodeAuth, opts)
		if !ok || err != nil {
			return nil, nil, err
		}

		token, err = a.CreateAccessTokenWith(&AccessTokenOptions{
			UserID:   opts.UserID,
			ClientID: opts.ClientID,
			Scope:    opts.Scope,
		})
		if err != nil {
			return nil, nil, err
		}

		refreshToken, err = a.CreateRefreshTokenWith(&RefreshTokenOptions{
			UserID:   opts.UserID,
			ClientID: opts.ClientID,
			Scope:    opts.Scope,
		})
		if err != nil {
			return nil, nil, err
		}

		return token, refreshToken, nil
	case ImplicitAuth:
		ok, err := a.AuthRefreshToken(ImplicitAuth, opts)
		if !ok || err != nil {
			return nil, nil, err
		}

		token, err = a.CreateAccessTokenWith(&AccessTokenOptions{
			UserID:   opts.UserID,
			ClientID: opts.ClientID,
			Scope:    opts.Scope,
		})
		if err != nil {
			return nil, nil, err
		}

		refreshToken, err = a.CreateRefreshTokenWith(&RefreshTokenOptions{
			UserID:   opts.UserID,
			ClientID: opts.ClientID,
			Scope:    opts.Scope,
		})
		if err != nil {
			return nil, nil, err
		}

		return token, refreshToken, nil
	case PasswordAuth:
		ok, err := a.AuthRefreshToken(PasswordAuth, opts)
		if !ok || err != nil {
			return nil, nil, err
		}

		token, err = a.CreateAccessTokenWith(&AccessTokenOptions{
			UserID:   opts.UserID,
			ClientID: opts.ClientID,
			Scope:    opts.Scope,
		})
		if err != nil {
			return nil, nil, err
		}

		refreshToken, err = a.CreateRefreshTokenWith(&RefreshTokenOptions{
			UserID:   opts.UserID,
			ClientID: opts.ClientID,
			Scope:    opts.Scope,
		})
		if err != nil {
			return nil, nil, err
		}

		return token, refreshToken, nil
	default:
		return nil, nil, errors.New("bad auth-type")
	}
}
