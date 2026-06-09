package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/imohashi/url-shortener-service/internal/infra/logger"
	"github.com/joho/godotenv"
)

const (
	serverPortEnv = "SERVER_PORT"
	defaultPort   = "8080"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	if err := router.SetTrustedProxies(nil); err != nil {
		logger.Error(err.Error())
	}

	port := getServerPort()

	logger.Info("Server listening on port: " + port)

	if err := gin.Default().Run(":" + port); err != nil {
		logger.Error(err.Error())
	}
}

func getServerPort() string {
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file")
	}

	port := os.Getenv(serverPortEnv)

	if port == "" {
		logger.Info("No port specified, using default port: " + defaultPort)
		return defaultPort
	}

	return port
}
