package main

import (
	"log"

	"github.com/JuanCJR/task-controller/internal/config"
	"github.com/JuanCJR/task-controller/internal/database"
	"github.com/JuanCJR/task-controller/internal/handler"
	"github.com/JuanCJR/task-controller/internal/repository"
	"github.com/JuanCJR/task-controller/internal/server"
	"github.com/JuanCJR/task-controller/internal/service"
)

// @title Task Controller API
// @version 1.0
// @description     API para gestión de tareas con RBAC
// @host            localhost:8081
// @BasePath        /api/v1/task-controller
func main() {
	//Load Env variables
	cfg := config.LoadConfig()

	//Db Connection
	db, err := database.NewConnection(cfg.DB)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Repository Creation
	userRepo := repository.NewUserRepository(db)

	//Service Creation
	userService := service.NewUserService(userRepo)

	//Handler Creation
	userHandler := handler.NewUserHandler(userService)

	handlers := []handler.RouteRegister{
		userHandler,
	}

	//Init server
	router := server.NewServer(cfg, handlers)

	log.Printf("Server listening on port %s", cfg.APP.Port)
	router.Run(":" + cfg.APP.Port)
}
