package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type AppConfig struct {
	Port                 string //HTPP server port
	DefaultAdminEmail    string //Default email for the initial admin user created at startup
	DefaultAdminPassword string //Default password for the initial admin user created at startup
	ExecuteSeed          bool   //Flag to determine if the database seeding should be executed at startup
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

func requiredEnv(key string, errors *[]string) (string, error) {
	val := os.Getenv(key)

	if val == "" {
		*errors = append(*errors, fmt.Sprintf("Environment variable missing %s", key))
	}
	return val, nil
}

func LoadConfig() *Config {

	log.Println("Loading env variables..")

	envErr := godotenv.Load()

	if envErr != nil {
		log.Printf("Error loading .env file: %v", envErr)
	}

	var errors []string

	var err error
	config.APP.Port, err = requiredEnv("APP_PORT", &errors)

	config.APP.DefaultAdminEmail, err = requiredEnv("DEFAULT_ADMIN_EMAIL", &errors)

	config.APP.DefaultAdminPassword, err = requiredEnv("DEFAULT_ADMIN_PASSWORD", &errors)

	executeSeedStr, err := requiredEnv("EXECUTE_SEED", &errors)

	executeSeed, err := strconv.ParseBool(executeSeedStr)

	if err != nil {
		errors = append(errors, "Environment variable EXECUTE_SEED must be a valid boolean (true/false)")
	}

	config.APP.ExecuteSeed = executeSeed

	// Auth
	config.Auth.JwtSecret, err = requiredEnv("JWT_SECRET", &errors)

	tokenExp, err := requiredEnv("TOKEN_EXPIRATION", &errors)

	tokenExpInt, err := strconv.Atoi(tokenExp)

	if err != nil {
		errors = append(errors, "Environment variable TOKEN_EXPIRATION must be a valid integer")
	}
	config.Auth.TokenExpiration = tokenExpInt

	// Database
	config.DB.Host, err = requiredEnv("DB_HOST", &errors)
	if err != nil {
		errors = append(errors, err.Error())
	}

	dbPort, err := requiredEnv("DB_PORT", &errors)
	if err != nil {
		errors = append(errors, err.Error())
	}
	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		errors = append(errors, "Environment variable DB_PORT must be a valid integer")
	}
	config.DB.Port = dbPortInt

	config.DB.User, err = requiredEnv("DB_USER", &errors)
	if err != nil {
		errors = append(errors, err.Error())
	}

	config.DB.Password, err = requiredEnv("DB_PASSWORD", &errors)
	if err != nil {
		errors = append(errors, err.Error())
	}

	config.DB.Name, err = requiredEnv("DB_NAME", &errors)
	if err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		log.Printf("Missing environment variables: %s", strings.Join(errors, "\n - "))
		os.Exit(1)
	}
	return &config
}

func GetConfig() *Config {
	return &config
}
