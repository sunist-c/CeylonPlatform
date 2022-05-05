package api

import (
	"CeylonPlatform/middleware/api"
	"github.com/gin-gonic/gin"
)

func init() {
	api.AddService(Service{})
}

type HealthServiceRequest struct {
}

type HealthServiceResponse struct {
}

type Service struct {
}

func (s Service) Handlers() []api.ServiceHandler {
	return []api.ServiceHandler{
		{
			Name:   "",
			Url:    "",
			Method: api.GET,
			Dependence: api.NewContextWith(&api.ContextOptions{
				DataBase: nil,
				Redis:    nil,
				Logger:   nil,
			}),
			Handler: func(req *gin.Context, ctx *api.Context) {

			},
			Router: nil,
		},
	}
}
