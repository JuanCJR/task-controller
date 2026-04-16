package config

import (
	"log"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type AppConfig struct {
	Port string //HTPP server port
}

type AuthConfig struct {
	JwtSecret       string //Secret key for JWT token generation and validation
	TokenExpiration int    //Token expiration time in minutes
}

type Config struct {
	APP  AppConfig      //Application configuration
	Auth AuthConfig     //Authentication configuration
	DB   DatabaseConfig //Database configuration
}

var config Config

func LoadConfig() {
	var errors []string

	config.APP.Port = os.Getenv("APP_PORT")
	if config.APP.Port == "" {
		errors = append(errors, "Enviroment variable missing APP_PORT")
		os.Exit(1)
	}

	// Auth
	config.Auth.JwtSecret = os.Getenv("JWT_SECRET")
	if config.Auth.JwtSecret == "" {
		errors = append(errors, "Environment variable missing JWT_SECRET")
	}

	tokenExp := os.Getenv("TOKEN_EXPIRATION")
	if tokenExp == "" {
		errors = append(errors, "Environment variable missing TOKEN_EXPIRATION")
	}
	tokenExpInt, err := strconv.Atoi(tokenExp)
	if err != nil {
		errors = append(errors, "Environment variable TOKEN_EXPIRATION must be a valid integer")
	}
	config.Auth.TokenExpiration = tokenExpInt

	// Database
	config.DB.Host = os.Getenv("DB_HOST")
	if config.DB.Host == "" {
		errors = append(errors, "Environment variable missing DB_HOST")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		errors = append(errors, "Environment variable missing DB_PORT")
	}
	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		errors = append(errors, "Environment variable DB_PORT must be a valid integer")
	}
	config.DB.Port = dbPortInt

	config.DB.User = os.Getenv("DB_USER")
	if config.DB.User == "" {
		errors = append(errors, "Environment variable missing DB_USER")
	}

	config.DB.Password = os.Getenv("DB_PASSWORD")
	if config.DB.Password == "" {
		errors = append(errors, "Environment variable missing DB_PASSWORD")
	}

	config.DB.Name = os.Getenv("DB_NAME")
	if config.DB.Name == "" {
		errors = append(errors, "Environment variable missing DB_NAME")
	}

	if len(errors) > 0 {
		for _, err := range errors {
			log.Println(err)
		}
		os.Exit(1)
	}

}

func GetConfig() *Config {
	return &config
}
