package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/kinoakter/openvpn-pki-go/internal/pki"
	"github.com/kinoakter/openvpn-pki-go/internal/repository"
	"github.com/kinoakter/openvpn-pki-go/internal/repository/entity"
	"github.com/kinoakter/openvpn-pki-go/log"
)

var CaService CAService = &caService{}

type CAService interface {
	CreateCA(context context.Context, serverName string, validYears int) error
}

type caService struct {
}

func (s *caService) CreateCA(context context.Context, serverName string, validYears int) error {
	certPEM, keyPEM, err := pki.CreateCACert(serverName, validYears)
	if err != nil {
		log.Errorf("failed to create CA: %v", err)
		return err
	}

	saveErr := repository.Ca.Save(context, &entity.CA{
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
