package server

import (
	"github.com/JuanCJR/task-controller/internal/config"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg *config.Config) *gin.Engine {
	router := gin.Default()
	return router
}
