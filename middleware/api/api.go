package api

import (
	"CeylonPlatform/middleware/logs"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"xorm.io/xorm"
)

var (
	serviceInterfaces []ServiceInterface = make([]ServiceInterface, 0, 16)
)

type ContextType string

const (
	Database ContextType = "db"
	Redis    ContextType = "redis"
	Logger   ContextType = "logger"
	// Authenticator ContextType = "authenticator"
)

type HttpMethod string

const (
	GET     HttpMethod = "GET"
	POST    HttpMethod = "POST"
	PATCH   HttpMethod = "PATCH"
	DELETE  HttpMethod = "DELETE"
	PUT     HttpMethod = "PUT"
	HEAD    HttpMethod = "HEAD"
	OPTIONS HttpMethod = "OPTIONS"
	ANY     HttpMethod = "ANY"
)

type Handler func(req *gin.Context, ctx *Context)

// Context 服务依赖环境
type Context struct {
	dictionary map[ContextType]interface{}
}

// ContextOptions 服务依赖环境选项
type ContextOptions struct {
	DataBase *xorm.Engine
	Redis    *redis.Client
	Logger   *logs.Logger
}

// NewContext 创建空的服务依赖环境
func NewContext() *Context {
	return &Context{dictionary: make(map[ContextType]interface{}, 4)}
}

// NewContextWith 使用自定义选项创建新的服务依赖环境
func NewContextWith(opts *ContextOptions) *Context {
	ctx := Context{dictionary: make(map[ContextType]interface{}, 4)}
	if opts != nil {
		if opts.DataBase != nil {
			ctx.dictionary[Database] = opts.DataBase
		}
		if opts.Redis != nil {
			ctx.dictionary[Redis] = opts.Redis
		}
		if opts.Logger != nil {
			ctx.dictionary[Logger] = opts.Logger
		}
	}
	return &ctx
}

// Get 获取服务依赖环境中的依赖
func (c Context) Get(key ContextType) (value interface{}, ok bool) {
	value, ok = c.dictionary[key]
	return
}

// ServiceHandler 服务处理器，用于处理某一类请求
type ServiceHandler struct {
	Name       string
	Uri        string
	Method     HttpMethod
	Handler    Handler
	Dependence *Context
	Router     *gin.RouterGroup
}

// ServiceInterface 服务接口，用于向公网提供服务
type ServiceInterface interface {
	Handlers() []ServiceHandler
}

// Bind 将各ServiceInterface的各Handler绑定到其对应的RouterGroup上
func Bind() {
	for _, serviceInterface := range serviceInterfaces {
		for _, handler := range serviceInterface.Handlers() {
			switch handler.Method {
			case GET:
				handler.Router.GET(handler.Uri, func(context *gin.Context) {
					handler.Handler(context, handler.Dependence)
				})
			case POST:
				handler.Router.POST(handler.Uri, func(context *gin.Context) {
					handler.Handler(context, handler.Dependence)
				})
			case PATCH:
				handler.Router.PATCH(handler.Uri, func(context *gin.Context) {
					handler.Handler(context, handler.Dependence)
				})
			case PUT:
				handler.Router.PUT(handler.Uri, func(context *gin.Context) {
					handler.Handler(context, handler.Dependence)
				})
			case DELETE:
				handler.Router.DELETE(handler.Uri, func(context *gin.Context) {
					handler.Handler(context, handler.Dependence)
				})
			case HEAD:
				handler.Router.HEAD(handler.Uri, func(context *gin.Context) {
					handler.Handler(context, handler.Dependence)
				})
			case OPTIONS:
				handler.Router.OPTIONS(handler.Uri, func(context *gin.Context) {
					handler.Handler(context, handler.Dependence)
				})
			case ANY:
				handler.Router.Any(handler.Uri, func(context *gin.Context) {
					handler.Handler(context, handler.Dependence)
				})
			default:
				handler.Router.Any(handler.Uri, func(context *gin.Context) {
					handler.Handler(context, handler.Dependence)
				})
			}
		}
	}
}

// AddService 向API中间件添加一个暴露于公网的服务
func AddService(serviceInterface ServiceInterface) {
	serviceInterfaces = append(serviceInterfaces, serviceInterface)
}
