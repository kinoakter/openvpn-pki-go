package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kinoakter/openvpn-pki-go/internal/config"
	"log"
)

var DB *pgxpool.Pool

func Connect(cfg *config.Config) {
	var err error
	DB, err = pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	fmt.Println("Connected to database")
}
