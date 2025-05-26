package entity

import "time"

type CA struct {
	CommonName  string
	Certificate string
	PrivateKey  string
	CreatedAt   time.Time
}
