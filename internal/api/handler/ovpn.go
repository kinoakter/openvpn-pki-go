package handler

import (
	"github.com/kinoakter/openvpn-pki-go/internal/domain/service"
	"github.com/labstack/echo/v4"
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

func (h *OVPNHandler) RegisterRoutes(router *echo.Group) {
	router.GET("/server-config/:server", h.GetServerOVPNConfig)
	router.GET("/client-config/:client_cn", h.GetClientOVPNConfig)
}

func (h *OVPNHandler) GetServerOVPNConfig(c echo.Context) error {
	serverName := c.Param("server")
	config, err := h.ovpnService.GenerateServerOVPNConfig(serverName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.String(http.StatusOK, config)
}

func (h *OVPNHandler) GetClientOVPNConfig(c echo.Context) error {
	//serverCN := c.Param("server_cn")
	clientCN := c.Param("client_cn")

	//if serverCN == "" {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "client CN is required"})
	//	return
	//}

	if clientCN == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "client CN is required"})
	}

	config, errCfg := h.ovpnService.GenerateClientOVPNConfig(clientCN)
	if errCfg != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": errCfg.Error()})
	}

	return c.String(http.StatusOK, config)
}
