package psql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kinoakter/openvpn-pki-go/log"
)

func MustConnect(context context.Context, dbURL string) (db *pgxpool.Pool) {
	var err error
	db, err = pgxpool.New(context, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return
	}

	log.Infof("Connected to database")

	return
}
