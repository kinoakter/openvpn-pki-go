package entity

import "time"

type ClientCert struct {
	CommonName          string
	ServerCommonName    string
	Certificate         string
	PrivateKey          string
	TlsCryptV2ClientKey string
	ExpiresAt           time.Time
	CreatedAt           time.Time
}
