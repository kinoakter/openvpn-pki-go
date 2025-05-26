package psql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kinoakter/openvpn-pki-go/log"
)

const schemaInit = `
	CREATE TABLE IF NOT EXISTS ca (
		common_name TEXT PRIMARY KEY,
		certificate TEXT NOT NULL,
		private_key TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS client_certificate (
		common_name TEXT PRIMARY KEY,
		server_common_name TEXT NOT NULL,
		certificate TEXT NOT NULL,
		private_key TEXT NOT NULL,
		tls_crypt_v2_key TEXT,
		revoked BOOLEAN DEFAULT FALSE,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMPTZ DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS server_certificate (
		common_name TEXT PRIMARY KEY,
		certificate TEXT NOT NULL,
		private_key TEXT NOT NULL,
		tls_crypt_v2_server_key TEXT,
		created_at TIMESTAMPTZ DEFAULT now()
	);
`

func Migrate(ctx context.Context, db *pgxpool.Pool) {
	_, err := db.Exec(ctx, schemaInit)

	if err != nil {
		log.Fatalf("failed to run db migrations: %v", err)
	}

	log.Infof("Database migration completed (if needed)")
}
