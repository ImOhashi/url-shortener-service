package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imohashi/url-shortener-service/docs"
)

func TestGetServerPortReturnsDefaultWhenEnvIsEmpty(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv(serverPortEnv, "")

	if got := getServerPort(); got != defaultPort {
		t.Fatalf("getServerPort() = %q, want %q", got, defaultPort)
	}
}

func TestGetServerPortReturnsEnvValue(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv(serverPortEnv, "9090")

	if got := getServerPort(); got != "9090" {
		t.Fatalf("getServerPort() = %q, want %q", got, "9090")
	}
}

func TestConfigureSwagger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalBasePath := docs.SwaggerInfo.BasePath
	originalHost := docs.SwaggerInfo.Host
	t.Cleanup(func() {
		docs.SwaggerInfo.BasePath = originalBasePath
		docs.SwaggerInfo.Host = originalHost
	})

	router := gin.New()

	configureSwagger(router, "9090")

	if docs.SwaggerInfo.BasePath != "/" {
		t.Fatalf("docs.SwaggerInfo.BasePath = %q, want %q", docs.SwaggerInfo.BasePath, "/")
	}

	if docs.SwaggerInfo.Host != "localhost:9090" {
		t.Fatalf("docs.SwaggerInfo.Host = %q, want %q", docs.SwaggerInfo.Host, "localhost:9090")
	}

	for _, route := range router.Routes() {
		if route.Method == "GET" && route.Path == "/swagger/*any" {
			return
		}
	}

	t.Fatal("swagger route GET /swagger/*any was not registered")
}
