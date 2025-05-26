package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

const DefaultClientCertValidityDays = 1

// IssueClientCertificate generates a new client certificate signed by the provided CA.
// It creates an ECDSA keypair and issues a certificate valid for client authentication.
//
// Parameters:
//   - caCertPEM: PEM-encoded CA certificate used to sign the client certificate
//   - caKeyPEM: PEM-encoded CA private key used for signing
//   - tlsCryptV2ServerKey: Server's TLS-Crypt-v2 key used to generate client specific key
//   - clientCommonName: Common Name (CN) to be used in the client certificate
//   - expiresAt: Time when the certificate will expire
//
// Returns:
//   - certPEM: PEM-encoded client certificate
//   - keyPEM: PEM-encoded client private key
//   - tlsCryptV2ClientKey: Generated client-specific TLS-Crypt-v2 key
//   - err: Error if any step of the certificate generation process fails
//
// The function performs the following steps:
//  1. Parses the CA certificate and private key
//  2. Generates a new ECDSA keypair for the client
//  3. Creates and signs the client certificate using the CA
//  4. Generates a client-specific TLS-Crypt-v2 key
func IssueClientCertificate(
	caCertPEM,
	caKeyPEM,
	tlsCryptV2ServerKey,
	clientCommonName string,
	expiresAt time.Time,
) (certPEM, keyPEM []byte, tlsCryptV2ClientKey string, err error) {

	// Decode CA certificate
	caCert, caParseErr := parseCertificate(caCertPEM)
	if caParseErr != nil {
		err = fmt.Errorf("failed to parse CA cert: %v", caParseErr)
		return
	}

	// Decode CA private key
	caKey, caKeyParseErr := parseECPrivateKey(caKeyPEM)
	if caKeyParseErr != nil {
		err = fmt.Errorf("failed to parse CA key: %v", caKeyParseErr)
		return
	}

	// 2. Generate the new ECDSA keypair for the client cert
	clientKey, cliKeyErr := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if cliKeyErr != nil {
		err = fmt.Errorf("failed to generate client key: %v", cliKeyErr)
		return
	}

	// 3. Build client certificate template
	serial := sha512.Sum512([]byte(uuid.New().String())) // Generate a unique serial number from UUID
	tmpl := &x509.Certificate{
		SerialNumber: new(big.Int).SetBytes(serial[:20]),
		Subject: pkix.Name{
			CommonName: clientCommonName,
		},
		NotBefore: time.Now().UTC(),
		NotAfter:  expiresAt.UTC(),
		KeyUsage:  x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	// Generate Subject Key Identifier (SKI) from client's public key
	spki, spkiErr := x509.MarshalPKIXPublicKey(&clientKey.PublicKey)
	if spkiErr != nil {
		err = fmt.Errorf("failed to marshal public key: %v", spkiErr)
		return
	}
	ski := sha512.Sum512_256(spki)
	tmpl.SubjectKeyId = ski[:]

	// Set Authority Key Identifier (AKI) from CA cert if present
	if len(caCert.SubjectKeyId) > 0 {
		tmpl.AuthorityKeyId = caCert.SubjectKeyId
	}

	// 4. Sign the client certificate with the CA
	der, derErr := x509.CreateCertificate(rand.Reader, tmpl, caCert, &clientKey.PublicKey, caKey)
	if derErr != nil {
		err = fmt.Errorf("failed to sign client cert: %v", derErr)
		return
	}

	// Encode signed certificate and private key to PEM format
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})

	privBytes, privErr := x509.MarshalECPrivateKey(clientKey)
	if privErr != nil {
		err = fmt.Errorf("failed to marshal client private key: %v", privErr)
		return
	}

	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})

	// 5. Generate tls-crypt-v2 client key
	tlsCryptV2ClientKey, err = GenerateTlsCryptV2ClientKey(tlsCryptV2ServerKey, clientCommonName)
	if err != nil {
		err = fmt.Errorf("failed to generate tls-crypt-v2 envelope: %v", err)
		return
	}

	return
}
