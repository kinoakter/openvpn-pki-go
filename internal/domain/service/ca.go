package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/entity"
	"github.com/kinoakter/openvpn-pki-go/internal/pki"
	"github.com/kinoakter/openvpn-pki-go/log"
)

type CARepository interface {
	LoadByServerName(context context.Context, serverName string) (*entity.CA, error)
	Save(context.Context, *entity.CA) (*entity.CA, error)
}

type CAService struct {
	ctx    context.Context
	caRepo CARepository
}

func NewCAService(ctx context.Context, caRepo CARepository) *CAService {
	return &CAService{
		ctx:    ctx,
		caRepo: caRepo,
	}
}

func (s *CAService) CreateCA(serverName string, validYears int) error {
	certPEM, keyPEM, err := pki.CreateCACert(serverName, validYears)
	if err != nil {
		log.Errorf("failed to create CA: %v", err)
		return err
	}

	_, saveErr := s.caRepo.Save(s.ctx, &entity.CA{
		UUID:        uuid.New(),
		Name:        serverName,
		Certificate: string(certPEM),
		PrivateKey:  string(keyPEM),
	})

	if saveErr != nil {
		log.Errorf("failed to save CA: %v", err)
		return err
	}

	return nil
}
