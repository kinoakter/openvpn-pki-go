package service

import (
	"context"
	"fmt"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/entity"
	"github.com/kinoakter/openvpn-pki-go/internal/pki"
)

type ServerCertificateRepository interface {
	LoadByServerName(context context.Context, serverName string) (*entity.ServerCert, error)
	Save(context.Context, *entity.ServerCert) (*entity.ServerCert, error)
}

type ServerCertificateService struct {
	ctx          context.Context
	repository   ServerCertificateRepository
	caRepository CARepository
}

func NewServerCertificateService(
	ctx context.Context,
	repository ServerCertificateRepository,
	caRepository CARepository,
) *ServerCertificateService {
	return &ServerCertificateService{
		ctx:          ctx,
		repository:   repository,
		caRepository: caRepository,
	}
}

func (s *ServerCertificateService) IssueNewServerCert(serverName string) error {
	ca, err := s.caRepository.LoadByServerName(s.ctx, serverName)
	if err != nil {
		return fmt.Errorf("failed to load CA by server name %s: %v", serverName, err)
	}

	certPEM, keyPEM, tlsCryptV2ServerKey, createErr := pki.IssueServerCertificate(serverName, ca.Certificate, ca.PrivateKey, pki.DefaultServerCertValidityYears)
	if createErr != nil {
		return fmt.Errorf("failed to create server cert: %v", createErr)
	}

	serverCert := &entity.ServerCert{
		CommonName:          serverName,
		Certificate:         string(certPEM),
		PrivateKey:          string(keyPEM),
		TlsCryptV2ServerKey: tlsCryptV2ServerKey,
	}

	_, saveErr := s.repository.Save(s.ctx, serverCert)
	if saveErr != nil {
		return fmt.Errorf("failed to save server cert: %v", saveErr)
	}

	return nil
}
