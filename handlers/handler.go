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
	FallbackUrl string `json:"fallback"`
	Ttl         string  `json:"ttl"`
}

func CreateShortUrl(context *gin.Context) {
	creationRequest := UrlCreationRequest{
		OriginalUrl: "",
		Ttl: "NA",
	}
	if err := context.ShouldBindJSON(&creationRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nextId := store.GetNextId()
	shortUrl := shortener.GenerateShortLink(creationRequest.OriginalUrl, nextId)
	store.SaveUrlMapping(nextId, shortUrl, creationRequest.OriginalUrl, creationRequest.FallbackUrl, creationRequest.Ttl)

	context.JSON(201, gin.H{
		"message":   "Short URL created successfully",
		"short_url": "/" + shortUrl,
	})

}

func HandleShortUrlRedirect(context *gin.Context) {
	shortUrl := context.Param("shortUrl")
	initialUrl, err := store.RetrieveCompleteUrl(shortUrl)

	if err != nil {
		fallbackUrl, err := store.RetrieveFallbackUrl(shortUrl)

		if err != nil {
			context.Status(404)
		} else {
			context.Redirect(302, fallbackUrl)
		}

		return
	} else {
		context.Redirect(302, initialUrl)
	}


	headerBytes, _ := json.Marshal(context.Request.Header)
	store.UpdateLink(shortUrl, string(headerBytes), context.ClientIP())
}

func ReactivateShortUrl(context *gin.Context) {
	shortUrl := context.Param("shortUrl")
	ttl := context.DefaultQuery("ttl", "NA")
	err := store.ReactivateUrl(shortUrl, ttl)
	if err != nil {
		context.AbortWithError(400, err)
	}
}

func DeactivateShortUrl(context *gin.Context) {
	shortUrl := context.Param("shortUrl")
	store.DeactivateUrl(shortUrl)
}
