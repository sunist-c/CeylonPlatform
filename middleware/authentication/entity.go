package authentication

import "time"

// User 用户的结构，用于维护用户基本信息
type User struct {
	ID       string    `xorm:"pk varchar(32) unique notnull index"`
	Name     string    `xorm:"varchar(32) notnull"`
	Password string    `xorm:"varchar(32) notnull"`
	Scope    ScopeType `xorm:"varchar(32) notnull default('student')"`
	CreateAt time.Time `xorm:"notnull"`
	UpdateAt time.Time `xorm:"notnull"`
}

// Client client表的结构，用于维护client的基本信息
type Client struct {
	ID          string    `xorm:"notnull varchar(32) unique index pk"`
	Name        string    `xorm:"notnull varchar(32)"`
	Key         string    `xorm:"notnull varchar(32)"`
	Secret      string    `xorm:"notnull varchar(32)"`
	RedirectUri string    `xorm:"notnull varchar(255)"`
	Scope       ScopeType `xorm:"varchar(32) notnull default('student')"`
	Method      AuthType  `xorm:"varchar(32) notnull default('client')"`
	CreateAt    time.Time `xorm:"notnull"`
	UpdateAt    time.Time `xorm:"notnull"`
}

// AccessToken token表的结构，用于维护token
type AccessToken struct {
	ID       string    `xorm:"notnull varchar(32) unique index pk"`
	Token    string    `xorm:"notnull varchar(32)"`
	UserID   string    `xorm:"notnull varchar(32)"`
	ClientID string    `xorm:"notnull varchar(32)"`
	Scope    ScopeType `xorm:"notnull varchar(32) default('student')"`
	ExpireAt time.Time `xorm:"notnull"`
	CreateAt time.Time `xorm:"notnull"`
}

// RefreshToken refresh-token表的结构，用于维护refresh-token
type RefreshToken struct {
	ID       string    `xorm:"notnull varchar(32) unique index pk"`
	Token    string    `xorm:"notnull varchar(32)"`
	UserID   string    `xorm:"notnull varchar(32)"`
	ClientID string    `xorm:"notnull varchar(32)"`
	Scope    ScopeType `xorm:"notnull varchar(32) default('student')"`
	ExpireAt time.Time `xorm:"notnull"`
	CreateAt time.Time `xorm:"notnull"`
}

// AuthorizationCode code表的结构，用于维护授权码
type AuthorizationCode struct {
	ID          string    `xorm:"notnull varchar(32) unique index pk"`
	Code        string    `xorm:"notnull varchar(32)"`
	UserID      string    `xorm:"notnull varchar(32)"`
	ClientID    string    `xorm:"notnull varchar(32)"`
	Scope       ScopeType `xorm:"notnull varchar(32) default('student')"`
	ExpireAt    time.Time `xorm:"notnull"`
	CreateAt    time.Time `xorm:"notnull"`
	RedirectUri string    `xorm:"notnull varchar(255)"`
}

type ScopeType string
type AuthType string

// 定义四种OAuth模式
const (
	CodeAuth     AuthType = "code"
	ImplicitAuth AuthType = "implicit"
	ClientAuth   AuthType = "client"
	PasswordAuth AuthType = "password"
)

// 定义四种权限类型
const (
	Manager ScopeType = "manager"
	Teacher ScopeType = "teacher"
	Student ScopeType = "student"
	Guest   ScopeType = "guest"
)