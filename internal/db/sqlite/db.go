package sqlite

import (
	"database/sql"
	"github.com/kinoakter/openvpn-pki-go/log"
)

var DB *sql.DB

func Connect(path string) {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("failed to connect to sqlite: %v", err)
	}
	log.Infof("Connected to database")
}
