package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/service"
	"net/http"
)

type ServerCertHandler struct {
	serverCertService *service.ServerCertificateService
}

func NewServerCertHandler(serverCertService *service.ServerCertificateService) *ServerCertHandler {
	return &ServerCertHandler{
		serverCertService: serverCertService,
	}
}

func (h *ServerCertHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/server-cert", h.IssueNewServerCert)
}

func (h *ServerCertHandler) IssueNewServerCert(c *gin.Context) {
	type Req struct {
		ServerName string `json:"server_name" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.serverCertService.IssueNewServerCert(req.ServerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
