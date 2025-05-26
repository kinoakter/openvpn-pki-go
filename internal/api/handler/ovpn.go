package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/service"
	"net/http"
)

type OVPNHandler struct {
	ovpnService *service.OVPNService
}

func NewOVPNHandler(ovpnService *service.OVPNService) *OVPNHandler {
	return &OVPNHandler{
		ovpnService: ovpnService,
	}
}

func (h *OVPNHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/server-config/:server", h.GetServerOVPNConfig)
	router.GET("/client-config/:client_cn", h.GetClientOVPNConfig)
}

func (h *OVPNHandler) GetServerOVPNConfig(c *gin.Context) {
	serverName := c.Param("server")
	config, err := h.ovpnService.GenerateServerOVPNConfig(serverName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, config)
}

func (h *OVPNHandler) GetClientOVPNConfig(c *gin.Context) {
	//serverCN := c.Param("server_cn")
	clientCN := c.Param("client_cn")

	//if serverCN == "" {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "client CN is required"})
	//	return
	//}

	if clientCN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client CN is required"})
		return
	}

	config, errCfg := h.ovpnService.GenerateClientOVPNConfig(clientCN)
	if errCfg != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errCfg.Error()})
		return
	}
	c.String(http.StatusOK, config)
}
