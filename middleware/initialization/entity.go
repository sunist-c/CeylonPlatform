package initialization

import (
	"CeylonPlatform/middleware/logs"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gopkg.in/ini.v1"
	"xorm.io/xorm"
)

func init() {
	InitEntityList = make(map[string]ConfigCustomer)
	InitEntityList["redis"] = &RedisConfig{}
	InitEntityList["mysql"] = &MysqlConfig{}
	InitEntityList["service"] = &ServiceConfig{}
}

type ConfigCustomer interface {
	// Read 读取配置文件
	Read(configName string, file *ini.File) error

	// Init 实例化组件
	Init() error
}

var (
	InitEntityList map[string]ConfigCustomer
)

// 定义复用对象
var (
	DbConnection    *xorm.Engine     // 数据库连接池
	ApiRouter       *gin.RouterGroup // API路由
	Logger          *logs.Logger     // 日志器
	RedisConnection *redis.Client    // Redis连接池
	engine          *gin.Engine      // gin引擎
	startFunction   func() error     // 启动函数
	ServeMode       ServiceMode      // 服务状态
)

// InitEntities 初始化可复用实体
func InitEntities(filePath string) error {
	// 读取配置文件
	file, err := ini.Load(filePath)
	if err != nil {
		return err
	}

	// 实例化可服用实体
	for key, customer := range InitEntityList {
		if err = customer.Read(key, file); err != nil {
			return err
		}
		if err = customer.Init(); err != nil {
			return err
		}
	}

	// 检查实体状态
	if DbConnection == nil || ApiRouter == nil || Logger == nil || RedisConnection == nil {
		// todo: complete empty middlewares
		panic("todo: complete empty middlewares")
	}

	return nil
}
