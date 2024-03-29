package main

import (
	"fmt"
	store "github.com/HRKings/url-shortener/data"
	handler "github.com/HRKings/url-shortener/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil && os.Getenv("ENV_VARS_PROVIDED") != "true" {
		log.Fatal("Error loading .env file")
	}

	store.InitializeStore()

	server := gin.Default()

	// Enable Allow All CORS
	server.Use(cors.Default())

	trustedProxies, hasTrustedProxies := os.LookupEnv("TRUSTED_PROXIES")

	if hasTrustedProxies {
		log.Printf("Trusting proxies: %s", trustedProxies)
		server.SetTrustedProxies(strings.Split(trustedProxies, ";"))
	}

	server.POST("/", func(c *gin.Context) {
		handler.CreateShortUrl(c)
	})

	server.PUT("/:shortUrl", func(c *gin.Context) {
		handler.ReactivateShortUrl(c)
	})
	server.DELETE("/:shortUrl", func(c *gin.Context) {
		handler.DeactivateShortUrl(c)
	})

	server.GET("/:shortUrl", func(c *gin.Context) {
		handler.HandleShortUrlRedirect(c)
	})

	err = server.Run(":5000")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}
