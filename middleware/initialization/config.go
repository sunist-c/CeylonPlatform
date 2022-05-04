package initialization

import (
	"CeylonPlatform/middleware/logs"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"xorm.io/xorm"
)

// RedisConfig redis的配置结构
type RedisConfig struct {
	Address  string
	Port     string
	Password string
	DbName   int
}

func (r RedisConfig) Init() error {
	// 创建redis连接
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", r.Address, r.Port),
		Password: r.Password,
		DB:       r.DbName,
	})

	// 判断redis连接是否可用
	if _, err := client.Ping().Result(); err != nil {
		return err
	} else {
		RedisConnection = client
		return nil
	}
}

func (r *RedisConfig) Read(configName string, file *ini.File) error {
	r.Address = file.Section(configName).Key("address").String()
	r.Port = file.Section(configName).Key("port").String()
	r.Password = file.Section(configName).Key("password").String()
	dbName, err := file.Section(configName).Key("dbName").Int()

	// 校验数据是否合法
	if r.Address != "" && r.Port != "" && err == nil {
		r.DbName = dbName
		return nil
	} else {
		// todo: complete empty config field error - redis
		panic(errors.New("complete empty config field error"))
	}
}

// MysqlConfig mysql的配置结构
type MysqlConfig struct {
	Address           string
	Port              string
	Username          string
	Password          string
	DbName            string
	maxOpenConnection int
	maxIdleConnection int
}

func (m MysqlConfig) Init() error {
	// 创建mysql连接
	engine, err := xorm.NewEngine("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4",
		m.Username, m.Password, m.Address, m.Port, m.DbName))
	if err != nil {
		return err
	}

	// 检查连接可用性
	err = engine.Ping()
	if err != nil {
		return err
	} else {
		engine.SetMaxOpenConns(m.maxOpenConnection)
		engine.SetMaxIdleConns(m.maxIdleConnection)
		DbConnection = engine
		return nil
	}
}

func (m *MysqlConfig) Read(configName string, file *ini.File) error {
	m.Address = file.Section(configName).Key("address").String()
	m.Port = file.Section(configName).Key("port").String()
	m.Username = file.Section(configName).Key("username").String()
	m.Password = file.Section(configName).Key("password").String()
	m.DbName = file.Section(configName).Key("dbName").String()
	m.maxOpenConnection, _ = file.Section(configName).Key("maxOpenConnection").Int()
	m.maxIdleConnection, _ = file.Section(configName).Key("maxIdleConnection").Int()

	// 校验数据是否合法
	if m.Address != "" && m.Port != "" && m.Username != "" && m.Password != "" && m.DbName != "" {
		return nil
	} else {
		// todo: complete empty config field error - mysql
		panic(errors.New("complete empty config field error"))
	}
}

// ServiceConfig service的配置结构
type ServiceConfig struct {
	Mode    ServiceMode
	Port    string
	LogPath string
}

func (s ServiceConfig) Init() error {
	// 实例化gin
	e := gin.New()
	ApiRouter = e.Group("/")
	engine = e
	startFunction = s.run()

	// 实例化logger
	if outFile, err := os.OpenFile(fmt.Sprintf("%v/out.log", s.LogPath), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		return err
	} else {
		switch s.Mode {
		case Debug:
			Logger = logs.Default(outFile, "[DEBUG]", log.LstdFlags|log.Llongfile)
		case Develop:
			Logger = logs.Default(os.Stdout, "[DEVELOP]", log.Ltime|log.Lshortfile)
		case Product:
			Logger = logs.Default(outFile, "[INFO]", log.LstdFlags)
		default:
			Logger = logs.Default(outFile, "[INFO]", log.LstdFlags)
		}
	}

	return nil
}

func (s *ServiceConfig) Read(configName string, file *ini.File) error {
	mode := file.Section(configName).Key("mode").String()
	s.Port = file.Section(configName).Key("port").String()
	s.LogPath = file.Section(configName).Key("logPath").String()

	// 推断类型并赋默认值
	switch mode {
	case "debug":
		s.Mode = Debug
	case "product":
		s.Mode = Product
	case "develop":
		s.Mode = Develop
	default:
		s.Mode = Debug
	}
	if s.Port == "" {
		s.Port = "8080"
	}
	if s.LogPath == "" {
		s.LogPath = "log"
	}
	return nil
}

func (s ServiceConfig) run() func() error {
	return func() error {
		return engine.Run(fmt.Sprintf("0.0.0.0:%v", s.Port))
	}
}

type ServiceMode string

const (
	Debug   ServiceMode = "debug"
	Product ServiceMode = "product"
	Develop ServiceMode = "develop"
)
