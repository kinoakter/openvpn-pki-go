package handler

import (
	"github.com/kinoakter/openvpn-pki-go/internal/domain/service"
	"github.com/labstack/echo/v4"
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

func (h *ServerCertHandler) RegisterRoutes(router *echo.Group) {
	router.POST("/server-cert", h.IssueNewServerCert)
}

func (h *ServerCertHandler) IssueNewServerCert(c echo.Context) error {
	type Req struct {
		ServerName string `json:"server_name" binding:"required"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	err := h.serverCertService.IssueNewServerCert(req.ServerName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
