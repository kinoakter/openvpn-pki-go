package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/service"
	"net/http"
)

type ClientCertHandler struct {
	clientCertService *service.ClientCertificateService
}

func NewClientCertHandler(service *service.ClientCertificateService) *ClientCertHandler {
	return &ClientCertHandler{
		clientCertService: service,
	}
}

func (h *ClientCertHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/client-cert", h.IssueNewClientCert)
}

func (h *ClientCertHandler) IssueNewClientCert(c *gin.Context) {
	type Req struct {
		ServerName string `json:"server_name" binding:"required"`
		CommonName string `json:"common_name" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.clientCertService.IssueNewClientCert(req.ServerName, req.CommonName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
