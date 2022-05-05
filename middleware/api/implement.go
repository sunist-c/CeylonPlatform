package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Implement func(ctx *gin.Context)

// todo: implement default implement elegant
var (
	InternalServerError Implement = func(ctx *gin.Context) {
		ctx.String(http.StatusInternalServerError, "internal server error")
	}

	NotFoundError Implement = func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "page not found")
	}

	BadRequestError Implement = func(ctx *gin.Context) {
		ctx.String(http.StatusBadRequest, "bad request")
	}

	ForbiddenError Implement = func(ctx *gin.Context) {
		ctx.String(http.StatusForbidden, "forbidden")
	}
)
