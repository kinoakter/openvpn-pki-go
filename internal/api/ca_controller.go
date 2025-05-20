package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kinoakter/openvpn-pki-go/internal/service"
	"net/http"
)

type CAController struct {
	caService service.CAService
}

func NewCAController(caService service.CAService) *CAController {
	return &CAController{
		caService: caService,
	}
}

// CreateCA Creates a new CA for a given OpenVPN server
func (ctrl *CAController) CreateCA(c *gin.Context) {
	type Req struct {
		ServerName string `json:"server_name" binding:"required"`
		ValidYears int    `json:"valid_years" binding:"required"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.caService.CreateCA(context.TODO(), req.ServerName, req.ValidYears); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"status": "CA created for " + req.ServerName})
}
