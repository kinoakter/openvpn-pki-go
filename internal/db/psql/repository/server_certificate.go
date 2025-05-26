package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/entity"
)

type ServerCertificateRepository struct {
	db *pgxpool.Pool
}

func NewServerCertificateRepository(db *pgxpool.Pool) *ServerCertificateRepository {
	return &ServerCertificateRepository{
		db: db,
	}
}

func (r *ServerCertificateRepository) Save(ctx context.Context, cert *entity.ServerCert) (*entity.ServerCert, error) {
	_, err := r.db.Exec(ctx,
		`INSERT INTO server_certificate (common_name, certificate, private_key, tls_crypt_v2_server_key) VALUES ($1, $2, $3, $4)`,
		cert.CommonName,
		cert.Certificate,
		cert.PrivateKey,
		cert.TlsCryptV2ServerKey,
	)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func (r *ServerCertificateRepository) LoadByServerName(ctx context.Context, commonName string) (*entity.ServerCert, error) {
	var cert entity.ServerCert
	err := r.db.QueryRow(ctx,
		`SELECT common_name, certificate, private_key, tls_crypt_v2_server_key 
		FROM server_certificate 
		WHERE common_name = $1`,
		commonName).
		Scan(&cert.CommonName, &cert.Certificate, &cert.PrivateKey, &cert.TlsCryptV2ServerKey)

	if err != nil {
		return nil, err
	}

	return &cert, nil
}
