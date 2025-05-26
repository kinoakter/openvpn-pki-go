package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/entity"
)

type ClientCertificateRepository struct {
	db *pgxpool.Pool
}

func NewClientCertificateRepository(db *pgxpool.Pool) *ClientCertificateRepository {
	return &ClientCertificateRepository{
		db: db,
	}
}

func (r *ClientCertificateRepository) Save(ctx context.Context, clientCert *entity.ClientCert) (rv *entity.ClientCert, err error) {
	_, err = r.db.Exec(ctx, `
				INSERT INTO client_certificate (common_name, server_common_name, certificate, private_key, tls_crypt_v2_key, expires_at) 
				VALUES ($1, $2, $3, $4, $5, $6)
				`,
		clientCert.CommonName,
		clientCert.ServerCommonName,
		clientCert.Certificate,
		clientCert.PrivateKey,
		clientCert.TlsCryptV2ClientKey,
		clientCert.ExpiresAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to save client cert: %v", err)
	}
	rv = clientCert

	return
}

func (r *ClientCertificateRepository) LoadByCommonName(ctx context.Context, commonName string) (*entity.ClientCert, error) {
	var cliCert entity.ClientCert
	err := r.db.QueryRow(ctx, `
					SELECT common_name, 
					       server_common_name, 
					       certificate, 
					       private_key, 
					       tls_crypt_v2_key,
					       expires_at,
					       created_at 
					FROM client_certificate 
					WHERE common_name = $1 LIMIT 1
					`, commonName,
	).Scan(
		&cliCert.CommonName,
		&cliCert.ServerCommonName,
		&cliCert.Certificate,
		&cliCert.PrivateKey,
		&cliCert.TlsCryptV2ClientKey,
		&cliCert.ExpiresAt,
		&cliCert.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load client cert: %v", err)
	}

	return &cliCert, nil
}
