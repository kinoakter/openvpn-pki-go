package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kinoakter/openvpn-pki-go/internal/api"
	"github.com/kinoakter/openvpn-pki-go/internal/config"
	"github.com/kinoakter/openvpn-pki-go/internal/db"
	"log"
)

func main() {
	cfg := config.Load()
	db.Connect(cfg)

	router := gin.Default()
	api.SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
