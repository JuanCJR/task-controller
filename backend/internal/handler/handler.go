package handler

import "github.com/gin-gonic/gin"

type RouteRegister interface {
	RegisterRoutes(rg *gin.RouterGroup)
}
