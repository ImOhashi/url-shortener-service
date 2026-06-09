package main

import (
	"errors"
	"os"
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

func TestGetServerPortLoadsEnvFile(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	originalPort, hadOriginalPort := os.LookupEnv(serverPortEnv)
	if err := os.Unsetenv(serverPortEnv); err != nil {
		t.Fatalf("os.Unsetenv(%q) error = %v", serverPortEnv, err)
	}
	t.Cleanup(func() {
		if hadOriginalPort {
			_ = os.Setenv(serverPortEnv, originalPort)
			return
		}

		_ = os.Unsetenv(serverPortEnv)
	})

	if err := os.WriteFile(".env", []byte(serverPortEnv+"=7070\n"), 0o600); err != nil {
		t.Fatalf("os.WriteFile(.env) error = %v", err)
	}

	if got := getServerPort(); got != "7070" {
		t.Fatalf("getServerPort() = %q, want %q", got, "7070")
	}
}

func TestMainConfiguresAndRunsServer(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv(serverPortEnv, "6060")

	originalSetTrustedProxies := setTrustedProxies
	originalRunRouter := runRouter
	originalBasePath := docs.SwaggerInfo.BasePath
	originalHost := docs.SwaggerInfo.Host
	t.Cleanup(func() {
		setTrustedProxies = originalSetTrustedProxies
		runRouter = originalRunRouter
		docs.SwaggerInfo.BasePath = originalBasePath
		docs.SwaggerInfo.Host = originalHost
	})

	setTrustedProxies = func(_ *gin.Engine, trustedProxies []string) error {
		if trustedProxies != nil {
			t.Fatalf("trustedProxies = %v, want nil", trustedProxies)
		}

		return errors.New("trusted proxies error")
	}

	runRouter = func(router *gin.Engine, address ...string) error {
		if len(address) != 1 || address[0] != ":6060" {
			t.Fatalf("address = %v, want [:6060]", address)
		}

		if docs.SwaggerInfo.BasePath != "/" {
			t.Fatalf("docs.SwaggerInfo.BasePath = %q, want %q", docs.SwaggerInfo.BasePath, "/")
		}

		if docs.SwaggerInfo.Host != "localhost:6060" {
			t.Fatalf("docs.SwaggerInfo.Host = %q, want %q", docs.SwaggerInfo.Host, "localhost:6060")
		}

		for _, route := range router.Routes() {
			if route.Method == "GET" && route.Path == "/swagger/*any" {
				return errors.New("router error")
			}
		}

		t.Fatal("swagger route GET /swagger/*any was not registered")
		return nil
	}

	main()
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
