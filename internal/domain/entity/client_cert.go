package entity

import "github.com/google/uuid"

type ClientCert struct {
	UUID                uuid.UUID
	CommonName          string
	Certificate         string
	PrivateKey          string
	TlsCryptV2ClientKey string
	ServerName          string
}
