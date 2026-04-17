package main

import (
	"log"

	"github.com/JuanCJR/task-controller/internal/config"
	"github.com/JuanCJR/task-controller/internal/database"
	"github.com/JuanCJR/task-controller/internal/handler"
	"github.com/JuanCJR/task-controller/internal/middleware"
	"github.com/JuanCJR/task-controller/internal/model"
	"github.com/JuanCJR/task-controller/internal/repository"
	"github.com/JuanCJR/task-controller/internal/server"
	"github.com/JuanCJR/task-controller/internal/service"
	"github.com/gin-gonic/gin"
)

// @title Task Controller API
// @version 1.0
// @description     API para gestión de tareas con RBAC
// @host            localhost:8081
// @BasePath        /api/v1/task-controller
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Ingresa el token con el formato: Bearer {token}
func main() {
	//Load Env variables
	cfg := config.LoadConfig()

	//Db Connection
	db, err := database.NewConnection(cfg.DB)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//DB seeding
	if cfg.APP.ExecuteSeed {
		database.Seed(db, cfg.APP)
	}

	//Repository Creation
	userRepo := repository.NewUserRepository(db)
	permissionRepo :=
		repository.NewPermissionRepository(db)

	//Middleware Creation
	authMw := middleware.AuthMiddleware(cfg.Auth.JwtSecret)
	rbacMw := func(action model.Action, module model.Module) gin.HandlerFunc {
		return middleware.RBACMiddleware(permissionRepo, action, module)
	}

	//Service Creation
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, cfg.Auth)

	//Handler Creation
	userHandler := handler.NewUserHandler(userService, authMw, rbacMw)
	authHandler := handler.NewAuthHandler(authService, authMw)

	handlers := []handler.RouteRegister{
		userHandler,
		authHandler,
	}

	//Init server
	router := server.NewServer(cfg, handlers)

	log.Printf("Server listening on port %s", cfg.APP.Port)
	router.Run(":" + cfg.APP.Port)
}
