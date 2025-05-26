package main

import (
	"github.com/kinoakter/openvpn-pki-go/internal/app"
	"log"
)

func main() {
	a := app.NewApp()

	if err := a.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	a.WaitForExit()
}
