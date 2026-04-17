package server

import (
	_ "github.com/JuanCJR/task-controller/docs"
	"github.com/JuanCJR/task-controller/internal/config"
	"github.com/JuanCJR/task-controller/internal/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewServer(cfg *config.Config, handlers []handler.RouteRegister) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("api/v1/task-controller")

	for _, h := range handlers {
		h.RegisterRoutes(api)
	}

	return router
}
