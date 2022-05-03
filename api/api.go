package api

import "github.com/gin-gonic/gin"

type ContextType string

const (
	Database      ContextType = "db"
	Redis         ContextType = "redis"
	Logger        ContextType = "logger"
	Authenticator ContextType = "authenticator"
)

type Context struct {
	dictionary map[ContextType]interface{}
}

func NewContext() Context {
	return Context{
		dictionary: make(map[ContextType]interface{}, 4),
	}
}

func (c *Context) Add(key ContextType, value interface{}) (ok bool) {
	if _, ok = c.dictionary[key]; ok {
		return false
	} else {
		c.dictionary[key] = value
		return true
	}
}

func (c Context) Get(key ContextType) (value interface{}, ok bool) {
	value, ok = c.dictionary[key]
	return
}

type Handler func(req *gin.Context, ctx Context)
