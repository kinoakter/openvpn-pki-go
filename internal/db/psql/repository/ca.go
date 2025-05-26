package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/entity"
	"github.com/kinoakter/openvpn-pki-go/log"
)

type CaRepository struct {
	db *pgxpool.Pool
}

func NewCaRepository(db *pgxpool.Pool) *CaRepository {
	return &CaRepository{
		db: db,
	}
}

func (r *CaRepository) LoadByServerName(context context.Context, commonName string) (*entity.CA, error) {
	// 1. Load CA for the given server name
	var ca entity.CA
	err := r.db.QueryRow(context, `
				SELECT common_name, 
				       certificate, 
				       private_key 
				FROM ca 
				WHERE common_name = $1 LIMIT 1
				`, commonName,
	).Scan(&ca.CommonName, &ca.Certificate, &ca.PrivateKey)

	if err != nil {
		return nil, err
	}

	return &ca, nil
}

func (r *CaRepository) Save(context context.Context, ca *entity.CA) (rv *entity.CA, err error) {
	_, err = r.db.Exec(context,
		`INSERT INTO ca (common_name, certificate, private_key) VALUES ($1, $2, $3)`,
		ca.CommonName, ca.Certificate, ca.PrivateKey,
	)

	if err != nil {
		log.Errorf("failed to save CA: %v", err)
	}

	rv = ca

	return
}
