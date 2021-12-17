package main

import (
	"fmt"
	"log"

	store "github.com/HRKings/url-shortener/data"
	handler "github.com/HRKings/url-shortener/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	store.InitializeStore()

	server := gin.Default()

	server.POST("/", func(c *gin.Context) {
		handler.CreateShortUrl(c)
	})

	server.GET("/:shortUrl", func(c *gin.Context) {
		handler.HandleShortUrlRedirect(c)
	})

	err = server.Run(":5000")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}
