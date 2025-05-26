package handler

import (
	"github.com/kinoakter/openvpn-pki-go/internal/domain/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CAHandler struct {
	caService *service.CAService
}

func NewCAHandler(caService *service.CAService) *CAHandler {
	return &CAHandler{
		caService: caService,
	}
}

func (h *CAHandler) RegisterRoutes(router *echo.Group) {
	router.POST("/ca", h.CreateCA)
}

// CreateCA Creates a new CA for a given OpenVPN server
func (h *CAHandler) CreateCA(c echo.Context) error {
	type Req struct {
		ServerName string `json:"server_name" binding:"required"`
		ValidYears int    `json:"valid_years" binding:"required"`
	}

	var req Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.caService.CreateCA(req.ServerName, req.ValidYears); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "CA created for " + req.ServerName})
}
