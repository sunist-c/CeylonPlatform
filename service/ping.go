package service

import (
	"CeylonPlatform/middleware/api"
	"github.com/gin-gonic/gin"
)

func init() {
	api.AddService(Service{})
}

type Service struct {
}

func (s Service) Handlers() []*api.ServiceHandler {
	pingHandler := api.ServiceHandler{
		Name:   "ping Handler",
		Url:    "ping",
		Method: api.GET,
		Handler: func(req *gin.Context, ctx *api.Context) {
			req.String(200, "pong")
		},
		Dependence: nil,
		Router:     api.BaseRouter(),
	}
	return []*api.ServiceHandler{
		&pingHandler,
	}
}
