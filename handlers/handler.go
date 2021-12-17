package handler

import (
	"encoding/json"
	"net/http"

	store "github.com/HRKings/url-shortener/data"
	shortener "github.com/HRKings/url-shortener/utils"

	"github.com/gin-gonic/gin"
)

type UrlCreationRequest struct {
	OriginalUrl string `json:"url" binding:"required"`
}

func CreateShortUrl(context *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := context.ShouldBindJSON(&creationRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl := shortener.GenerateShortLink(creationRequest.OriginalUrl, store.GetNextId())
	store.SaveUrlMapping(shortUrl, creationRequest.OriginalUrl)

	context.JSON(201, gin.H{
		"message":   "Short URL created successfully",
		"short_url": "/" + shortUrl,
	})

}

func HandleShortUrlRedirect(context *gin.Context) {
	shortUrl := context.Param("shortUrl")
	initialUrl := store.RetrieveCompleteUrl(shortUrl)
	context.Redirect(302, initialUrl)

	headerBytes, _ := json.Marshal(context.Request.Header)
	store.UpdateLink(shortUrl, string(headerBytes), context.ClientIP())
}
