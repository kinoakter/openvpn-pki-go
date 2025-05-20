package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kinoakter/openvpn-pki-go/internal/db"
	"github.com/kinoakter/openvpn-pki-go/internal/repository/entity"
	"github.com/kinoakter/openvpn-pki-go/log"
)

var Ca CaRepository = &caRepository{
	db: db.DB,
}

type CaRepository interface {
	Save(context context.Context, ca *entity.CA) error
}

type caRepository struct {
	db *pgxpool.Pool
}

func (repo *caRepository) Save(context context.Context, ca *entity.CA) error {
	_, err := db.DB.Exec(context,
		`INSERT INTO ca (id, name, certificate, private_key) VALUES ($1, $2, $3, $4)`,
		ca.UUID, ca.Name, ca.Certificate, ca.PrivateKey,
	)

	if err != nil {
		log.Errorf("failed to save CA: %v", err)
		return err
	}

	return nil
}
