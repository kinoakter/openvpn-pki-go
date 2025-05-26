package entity

import "time"

type ServerCert struct {
	CommonName          string
	Certificate         string
	PrivateKey          string
	TlsCryptV2ServerKey string
	CreatedAt           time.Time
}
