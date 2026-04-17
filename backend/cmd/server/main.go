package main

import (
	"log"

	"github.com/JuanCJR/task-controller/internal/config"
	"github.com/JuanCJR/task-controller/internal/server"
)

func main() {
	cfg := config.LoadConfig()
	router := server.NewServer(cfg)

	log.Printf("Server listening on port %s", cfg.APP.Port)
	router.Run(":" + cfg.APP.Port)
}
