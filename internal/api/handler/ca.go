package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/service"
	"net/http"
)

type CAHandler struct {
	caService *service.CAService
	router    *gin.Engine
}

func NewCAHandler(caService *service.CAService) *CAHandler {
	return &CAHandler{
		caService: caService,
	}
}

func (h *CAHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/ca", h.CreateCA)
}

// CreateCA Creates a new CA for a given OpenVPN server
func (h *CAHandler) CreateCA(c *gin.Context) {
	type Req struct {
		ServerName string `json:"server_name" binding:"required"`
		ValidYears int    `json:"valid_years" binding:"required"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.caService.CreateCA(req.ServerName, req.ValidYears); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"status": "CA created for " + req.ServerName})
}
