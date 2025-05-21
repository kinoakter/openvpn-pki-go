package entity

import "github.com/google/uuid"

type ServerCert struct {
	UUID                uuid.UUID
	Name                string
	Certificate         string
	PrivateKey          string
	TlsCryptV2ServerKey string
}
