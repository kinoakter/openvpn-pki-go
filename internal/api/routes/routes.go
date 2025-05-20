package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kinoakter/openvpn-pki-go/internal/api"
)

func RegisterCaController(router *gin.Engine, caController *api.CAController) {
	router.POST("/ca", caController.CreateCA)
}
