package pki

import (
	"fmt"
	"github.com/kinoakter/openvpn-pki-go/internal/shell"
	"os"
	"path/filepath"
)

const (
	OpenVPN = "/opt/homebrew/sbin/openvpn"
)

// GenerateTlsCryptV2ServerKey generates a 1024-bit (128-byte) tls-crypt-v2 server key consisting of two 512-bit keys.
// It returns the raw server key blob (128 bytes) and the PEM-formatted full server key string.
func GenerateTlsCryptV2ServerKey() (tlsCryptV2ServerKey string, err error) {
	tlsCryptV2ServerKey, err = shell.ExecCommand(OpenVPN, "--genkey", "tls-crypt-v2-server")
	return
}

// GenerateTlsCryptV2ClientKey generates a tls-crypt-v2 client key using the provided server key.
// It returns the generated client key as a string or an error if the operation fails.
func GenerateTlsCryptV2ClientKey(serverKey string, clientName string) (tlsCryptV2ClientKey string, err error) {
	tmpFile, errCreate := os.CreateTemp("", fmt.Sprintf("tmp_tls-crypt-v2-client-%s.key", clientName))
	if errCreate != nil {
		return "", errCreate
	}
	defer os.Remove(tmpFile.Name())

	// Write the server key to the file
	_, errWrite := tmpFile.WriteString(serverKey)
	if errWrite != nil {
		return "", errWrite
	}
	defer tmpFile.Close()

	pathToKey, errPath := filepath.Abs(tmpFile.Name())
	if errPath != nil {
		return "", errPath
	}

	tlsCryptV2ClientKey, err = shell.ExecCommand(OpenVPN, "--tls-crypt-v2", pathToKey, "--genkey", "tls-crypt-v2-client")

	return
}
