package repository

import (
	"context"
	"github.com/kinoakter/openvpn-pki-go/internal/db"
	"github.com/kinoakter/openvpn-pki-go/internal/repository/entity"
	"github.com/kinoakter/openvpn-pki-go/log"
)

var ServerCert ServerCertRepository = &serverCertRepository{}

type ServerCertRepository interface {
	Save(ctx context.Context, serverCert *entity.ServerCert) error
}

type serverCertRepository struct {
	db db.Db
}

func (repo *serverCertRepository) Save(ctx context.Context, serverCert *entity.ServerCert) (err error) {
	_, err = repo.db.Exec(ctx,
		`INSERT INTO server_certificates (id, server_name, certificate, private_key, tls_crypt_v2_server_key) VALUES ($1, $2, $3, $4, $5)`,
		serverCert.UUID, serverCert.Name, serverCert.Certificate, serverCert.PrivateKey, serverCert.TlsCryptV2ServerKey,
	)

	if err != nil {
		log.Errorf("failed to save server cert: %v", err)
	}

	return
}
