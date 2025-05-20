package entity

import "github.com/google/uuid"

type CA struct {
	UUID        uuid.UUID
	Name        string
	Certificate string
	PrivateKey  string
}
