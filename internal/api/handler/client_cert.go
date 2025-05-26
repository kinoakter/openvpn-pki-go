package handler

import (
	"github.com/kinoakter/openvpn-pki-go/internal/api/mapper"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/service"
	"github.com/labstack/echo/v4"
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

func (h *ClientCertHandler) RegisterRoutes(router *echo.Group) {
	router.POST("/client-cert", h.IssueNewClientCert)
	router.GET("/client/:commonName", h.GetClientCert)
}

func (h *ClientCertHandler) IssueNewClientCert(c echo.Context) error {
	type Req struct {
		ServerName string `json:"server_name" binding:"required"`
		CommonName string `json:"common_name" binding:"required"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.clientCertService.IssueNewClientCert(req.ServerName, req.CommonName); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (h *ClientCertHandler) GetClientCert(c echo.Context) error {
	commonName := c.Param("commonName")
	cert, ca, err := h.clientCertService.GetClientCertMaterials(commonName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	response := mapper.ToClientCertsResponse(cert, ca)
	return c.JSON(http.StatusOK, response)
}
