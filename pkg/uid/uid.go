package uid

import (
	"crypto"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

type UidType string

const (
	Client       UidType = "client"
	User         UidType = "user"
	Token        UidType = "token"
	RefreshToken UidType = "refresh-token"
	AuthCode     UidType = "auth-code"
	ClientKey    UidType = "client-key"
	ClientSecret UidType = "client-secret"

	clientRandomLower       int = 1000000
	userRandomLower         int = 2000000
	tokenRandomLower        int = 3000000
	refreshTokenRandomLower int = 4000000
	authCodeRandomLower     int = 5000000
	clientKeyRandomLower    int = 6000000
	clientSecretRandomLower int = 7000000
)

// GenerateMd5Len16 生成16位长字符串的md5摘要
func GenerateMd5Len16(str, salt string) (md5 string) {
	return GenerateMd5Len32(str, salt)[8:24]
}

// GenerateMd5Len32 生成16位长字符串的md5摘要
func GenerateMd5Len32(str, salt string) (md5 string) {
	h := crypto.MD5.New()
	h.Write([]byte(str + salt))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateUid(uidType UidType) (uid string) {
	randNum := rand.Int() % 1000000
	switch uidType {
	case Client:
		randNum += clientRandomLower
	case User:
		randNum += userRandomLower
	case Token:
		randNum += tokenRandomLower
	case RefreshToken:
		randNum += refreshTokenRandomLower
	case AuthCode:
		randNum += authCodeRandomLower
	case ClientKey:
		randNum += clientKeyRandomLower
	case ClientSecret:
		randNum += clientSecretRandomLower
	default:
		randNum += 0
	}

	nanoTime := time.Now().UnixNano()
	unixTime := time.Now().Unix()

	uid = GenerateMd5Len16(strconv.Itoa(randNum), "") + GenerateMd5Len16(strconv.FormatInt(nanoTime, 10), strconv.FormatInt(unixTime, 10))
	return
}
