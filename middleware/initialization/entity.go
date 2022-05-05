package initialization

import (
	"CeylonPlatform/middleware/logs"
	"github.com/go-redis/redis"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"xorm.io/xorm"
)

func init() {
	initEntityList = make(map[string]ConfigCustomer)
	initEntityList["redis"] = &RedisConfig{}
	initEntityList["mysql"] = &MysqlConfig{}
	initEntityList["service"] = &ServiceConfig{}
}

type ConfigCustomer interface {
	// Read 读取配置文件
	Read(configName string, file *ini.File) error

	// Init 实例化组件
	Init() error
}

var (
	initEntityList map[string]ConfigCustomer
)

// 定义复用对象
var (
	DbConnection    *xorm.Engine  // 数据库连接池
	Logger          *logs.Logger  // 日志器
	RedisConnection *redis.Client // Redis连接池
	starFunction    func() error  // 启动函数
)

// InitEntities 初始化可复用实体
func InitEntities(filePath string) error {
	// 读取配置文件
	file, err := ini.Load(filePath)
	if err != nil {
		return err
	}

	// 实例化可服用实体
	for key, customer := range initEntityList {
		if err = customer.Read(key, file); err != nil {
			return err
		}
		if err = customer.Init(); err != nil {
			return err
		}
	}

	// 检查实体状态
	if DbConnection == nil || Logger == nil || RedisConnection == nil {
		// todo: complete empty middlewares
		panic("todo: complete empty middlewares")
	}

	return nil
}

func Close() (code int) {
	log.Println("closing entities")
	err := DbConnection.Close()
	if err != nil {
		log.Println("database connection close failed")
		return 1
	}
	err = RedisConnection.Close()
	if err != nil {
		log.Println("redis connection close failed")
		return 1
	}

	log.Println("entities closed")
	return 0
}

func ListenSignal(sigs chan os.Signal) {
	sig := <-sigs
	log.Println(sig)
	os.Exit(Close())
}
