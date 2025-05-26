package mapper

import (
	"github.com/kinoakter/openvpn-pki-go/internal/domain/entity"
	"github.com/kinoakter/openvpn-pki-go/pkg/api/dto"
)

func ToClientCertsResponse(cert *entity.ClientCert, ca *entity.CA) *dto.ClientCertResponse {
	if cert == nil {
		return nil
	}

	return &dto.ClientCertResponse{
		CaCrt:               ca.Certificate,
		ClientCrt:           cert.Certificate,
		ClientKey:           cert.PrivateKey,
		TlsCryptClientV2Key: cert.TlsCryptV2ClientKey,
		ExpiresAt:           cert.ExpiresAt.UnixMilli(),
	}
}
