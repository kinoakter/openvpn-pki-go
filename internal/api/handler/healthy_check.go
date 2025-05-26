package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterHealthyCheck(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
