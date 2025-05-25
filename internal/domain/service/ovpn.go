package service

import (
	"context"
	"fmt"
)

const ServerConfigTemplate = `
# port 2443
proto tcp
dev tun_udp0
explicit-exit-notify 1

server 172.20.0.0 255.255.240.0

push "redirect-gateway def1 bypass-dhcp"
push "dhcp-option DNS 172.30.0.1"

duplicate-cn
keepalive 15 60
tls-server
tls-timeout 5
topology subnet
fast-io
float
mute-replay-warnings
# txqueuelen 5000 # Not supported on  Mac OS

allow-compression no
max-clients 2000
persist-key
persist-tun
status status-udp0.log 4
status-version 2

verb 3
data-ciphers AES-256-GCM:AES-128-GCM:CHACHA20-POLY1305
data-ciphers-fallback AES-256-GCM
auth none
dh none

`

const ClientConfigTemplate = `
client
remote localhost 1194
proto tcp
dev tun
persist-tun
nobind

connect-retry 2
connect-retry-max 3
tls-client
remote-cert-tls server
resolv-retry infinite
allow-compression no
ifconfig-nowarn

verb 3
auth-nocache
preresolve

# Encryption
data-ciphers CHACHA20-POLY1305
data-ciphers-fallback AES-256-GCM
auth none

`

type OVPNService struct {
	ctx            context.Context
	caRepository   CARepository
	serverRepo     ServerCertificateRepository
	clientCertRepo ClientCertificateRepository
}

func NewOVPNService(
	ctx context.Context,
	caRepo CARepository,
	serverRepo ServerCertificateRepository,
	clientCertRepo ClientCertificateRepository) *OVPNService {
	return &OVPNService{
		ctx:            ctx,
		caRepository:   caRepo,
		serverRepo:     serverRepo,
		clientCertRepo: clientCertRepo,
	}
}

func (s *OVPNService) GenerateServerOVPNConfig(serverName string) (string, error) {
	// Load CA certificate
	ca, caErr := s.caRepository.LoadByServerName(s.ctx, serverName)
	if caErr != nil {
		return "", fmt.Errorf("failed to load CA by server name %s: %v", serverName, caErr)
	}

	// Load server cert entity
	serverCert, serverCertErr := s.serverRepo.LoadByServerName(s.ctx, serverName)
	if serverCertErr != nil {
		return "", fmt.Errorf("failed to load server cert by server name %s: %v", serverName, serverCertErr)
	}

	// Format the OpenVPN config string with embedded blocks
	return fmt.Sprintf(`# OpenVPN Server Configuration
%s
<cert>
%s</cert>

<key>
%s</key>

<ca>
%s</ca>

<tls-crypt-v2>
%s</tls-crypt-v2>
`, ServerConfigTemplate,
		serverCert.Certificate,
		serverCert.PrivateKey,
		ca.Certificate,
		serverCert.TlsCryptV2ServerKey,
	), nil

}

func (s *OVPNService) GenerateClientOVPNConfig(commonName string) (string, error) {
	cliCert, err := s.clientCertRepo.LoadByCommonName(s.ctx, commonName)
	if err != nil {
		return "", fmt.Errorf("failed to load client cert by common name %s: %v", commonName, err)
	}

	ca, caErr := s.caRepository.LoadByServerName(s.ctx, cliCert.ServerName)
	if caErr != nil {
		return "", fmt.Errorf("failed to load CA by server name %s: %v", cliCert.ServerName, caErr)
	}

	// Format the OpenVPN config string with embedded blocks
	clientConfig := fmt.Sprintf(`
%s
<cert>
%s</cert>

<key>
%s</key>

<ca>
%s</ca>

<tls-crypt-v2>
%s</tls-crypt-v2>
`, ClientConfigTemplate,
		cliCert.Certificate,
		cliCert.PrivateKey,
		ca.Certificate,
		cliCert.TlsCryptV2ClientKey)

	return clientConfig, nil
}
