package pki

import (
	"fmt"
	"github.com/kinoakter/openvpn-pki-go/internal/shell"
	"os"
	"path/filepath"
)

const (
	// OpenVPN specifies the path to the OpenVPN executable used for key generation
	OpenVPN = "/opt/homebrew/sbin/openvpn"
)

// GenerateTlsCryptV2ServerKey generates a new TLS-Crypt-v2 server key using OpenVPN's
// key generation functionality.This key is used to create client-specific TLS-Crypt-v2 keys.
//
// Returns:
//   - tlsCryptV2ServerKey: Generated server key as a string
//   - err: Error if key generation fails
func GenerateTlsCryptV2ServerKey() (tlsCryptV2ServerKey string, err error) {
	tlsCryptV2ServerKey, err = shell.ExecCommand(OpenVPN, "--genkey", "tls-crypt-v2-server")
	return
}

// GenerateTlsCryptV2ClientKey generates a client-specific TLS-Crypt-v2 key using
// the server's TLS-Crypt-v2 key. This creates a unique key for each client that
// can be used for additional authentication and encryption.
//
// Parameters:
//   - serverKey: The server's TLS-Crypt-v2 key used as base for generation
//   - clientName: Name of the client for which the key is being generated
//
// Returns:
//   - tlsCryptV2ClientKey: Generated client-specific key as a string
//   - err: Error if key generation fails
func GenerateTlsCryptV2ClientKey(serverKey string, clientName string) (tlsCryptV2ClientKey string, err error) {
	tmpFile, errCreate := os.CreateTemp("", fmt.Sprintf("tmp_tls-crypt-v2-client-%s.key", clientName))
	if errCreate != nil {
		return "", errCreate
	}
	defer func(name string) {
		_ = os.Remove(name)
	}(tmpFile.Name())

	// Write the server key to the file
	_, errWrite := tmpFile.WriteString(serverKey)
	if errWrite != nil {
		return "", errWrite
	}
	defer func(tmpFile *os.File) {
		_ = tmpFile.Close()
	}(tmpFile)

	pathToKey, errPath := filepath.Abs(tmpFile.Name())
	if errPath != nil {
		return "", errPath
	}

	tlsCryptV2ClientKey, err = shell.ExecCommand(OpenVPN, "--tls-crypt-v2", pathToKey, "--genkey", "tls-crypt-v2-client")

	return
}
