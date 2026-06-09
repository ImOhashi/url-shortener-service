package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/imohashi/url-shortener-service/docs"
	"github.com/imohashi/url-shortener-service/internal/infra/logger"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	serverPortEnv = "SERVER_PORT"
	defaultPort   = "8080"
)

// @title URL Shortener Service API
// @version 1.0
// @description API for shortening URLs.
// @host localhost:8000
// @BasePath /
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

	configureSwagger(router, port)

	if err := router.Run(":" + port); err != nil {
		logger.Error(err.Error())
	}
}

func configureSwagger(router *gin.Engine, port string) {
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = "localhost:" + port

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
