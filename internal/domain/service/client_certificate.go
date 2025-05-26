package service

import (
	"context"
	"fmt"
	"github.com/kinoakter/openvpn-pki-go/internal/domain/entity"
	"github.com/kinoakter/openvpn-pki-go/internal/pki"
	"time"
)

type ClientCertificateRepository interface {
	LoadByCommonName(ctx context.Context, commonName string) (*entity.ClientCert, error)
	Save(context.Context, *entity.ClientCert) (*entity.ClientCert, error)
}

type ClientCertificateService struct {
	ctx          context.Context
	repository   ClientCertificateRepository
	caRepository CARepository
	serverRepo   ServerCertificateRepository
}

func NewClientCertificateService(
	ctx context.Context,
	repository ClientCertificateRepository,
	caRepository CARepository,
	serverRepo ServerCertificateRepository) *ClientCertificateService {
	return &ClientCertificateService{
		ctx:          ctx,
		repository:   repository,
		caRepository: caRepository,
		serverRepo:   serverRepo,
	}
}

func (s *ClientCertificateService) IssueNewClientCert(serverName, clientCommonName string) error {
	ca, err := s.caRepository.LoadByServerName(s.ctx, serverName)
	if err != nil {
		return fmt.Errorf("failed to load CA by server name %s: %v", serverName, err)
	}

	serverCert, srvLoadErr := s.serverRepo.LoadByServerName(s.ctx, serverName)
	if srvLoadErr != nil {
		return fmt.Errorf("failed to load server cert by server name %s: %v", serverName, srvLoadErr)
	}

	expiresAt := time.Now().AddDate(0, 0, pki.DefaultClientCertValidityDays).UTC()
	certPEM, keyPEM, tlsCryptV2ClientKey, cliCertErr := pki.IssueClientCertificate(
		ca.Certificate,
		ca.PrivateKey,
		serverCert.TlsCryptV2ServerKey,
		clientCommonName,
		expiresAt,
	)
	if cliCertErr != nil {
		return fmt.Errorf("failed to create client cert: %v", cliCertErr)
	}

	clientCert := &entity.ClientCert{
		CommonName:          clientCommonName,
		ServerCommonName:    serverName,
		Certificate:         string(certPEM),
		PrivateKey:          string(keyPEM),
		TlsCryptV2ClientKey: tlsCryptV2ClientKey,
		ExpiresAt:           expiresAt,
	}

	_, saveErr := s.repository.Save(s.ctx, clientCert)
	if saveErr != nil {
		return fmt.Errorf("failed to save client cert: %v", saveErr)
	}

	return nil
}

func (s *ClientCertificateService) GetClientCertMaterials(commonName string) (*entity.ClientCert, *entity.CA, error) {
	cliCert, err := s.repository.LoadByCommonName(s.ctx, commonName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load client cert by common name %s: %v", commonName, err)
	}

	ca, caErr := s.caRepository.LoadByServerName(s.ctx, cliCert.ServerCommonName)
	if caErr != nil {
		return nil, nil, fmt.Errorf("failed to load CA by server name %s: %v", cliCert.ServerCommonName, caErr)
	}

	return cliCert, ca, nil
}
