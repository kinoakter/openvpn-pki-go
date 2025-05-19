// File: internal/pki/tlscryptv2_test.go
package pki

import (
	"log"
	"testing"
)

func TestGenerateTlsCryptV2ServerKey(t *testing.T) {
	key, err := GenerateTlsCryptV2ServerKey()
	if err != nil {
		t.Fatalf("GenerateTlsCryptV2ServerKey failed: %v", err)
	}

	log.Printf("Server key:\n%s", key)
}

func TestGenerateTlsCryptV2ClientKey(t *testing.T) {
	serverKey, err := GenerateTlsCryptV2ServerKey()
	if err != nil {
		t.Fatalf("GenerateTlsCryptV2ServerKey failed: %v", err)
	}

	gotTlsCryptV2ClientKey, clientKeyErr := GenerateTlsCryptV2ClientKey(serverKey, "client-0")
	if clientKeyErr != nil {
		t.Fatalf("GenerateTlsCryptV2ClientKey failed: %v", clientKeyErr)
	}

	log.Printf("Client key:\n%s", gotTlsCryptV2ClientKey)
}
