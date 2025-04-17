package db

import (
	"context"
	"log"
)

func Migrate() {
	ctx := context.Background()

	_, err := DB.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS ca (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		certificate TEXT NOT NULL,
		private_key TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS client_certificates (
		id UUID PRIMARY KEY,
		common_name TEXT NOT NULL,
		certificate TEXT NOT NULL,
		private_key TEXT NOT NULL,
		tls_crypt_v2_key TEXT,
		revoked BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT now()
	);
	`)

	if err != nil {
		log.Fatalf("failed to run DB migrations: %v", err)
	}

	log.Println("Database migration completed (if needed)")
}
