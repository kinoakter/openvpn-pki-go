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

const DefaultServerCertValidityYears = 2

// IssueServerCertificate generates a new server certificate signed by the provided CA.
// It creates an ECDSA keypair and issues a certificate valid for server authentication.
//
// Parameters:
//   - serverName: Name of the server to be used in the certificate's Common Name
//   - caCertPEM: PEM-encoded CA certificate used to sign the server certificate
//   - caKeyPEM: PEM-encoded CA private key used for signing
//   - validYears: Number of years the certificate will be valid for (defaults to DefaultServerCertValidityYears if 0)
//
// Returns:
//   - certPEM: PEM-encoded server certificate
//   - keyPEM: PEM-encoded server private key
//   - tlsCryptV2ServerKey: Generated TLS-Crypt-v2 server key
//   - err: Error if any step of the certificate generation process fails
//
// The function performs the following steps:
//  1. Parses the CA certificate and private key
//  2. Generates a new ECDSA keypair for the server
//  3. Creates and signs the server certificate using the CA
//  4. Generates a TLS-Crypt-v2 server key
func IssueServerCertificate(serverName, caCertPEM, caKeyPEM string, validYears int) (certPEM, keyPEM []byte, tlsCryptV2ServerKey string, err error) {
	if validYears == 0 {
		validYears = DefaultServerCertValidityYears // By default, 2 years.
	}

	// Decode CA certificate
	caCert, certErr := parseCertificate(caCertPEM)
	if certErr != nil {
		err = fmt.Errorf("failed to parse CA cert: %v", certErr)
		return
	}

	// Decode CA private key
	caKey, keyErr := parseECPrivateKey(caKeyPEM)
	if keyErr != nil {
		err = fmt.Errorf("failed to parse CA private key: %v", keyErr)
		return
	}

	// 2. Generate new keypair for the server cert
	serverKey, genErr := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if genErr != nil {
		err = fmt.Errorf("failed to generate server key: %v", genErr)
		return
	}

	// 3. Create server certificate template
	serial := sha512.Sum512([]byte(uuid.New().String()))
	sPki, pkiErr := x509.MarshalPKIXPublicKey(&serverKey.PublicKey)
	if pkiErr != nil {
		err = fmt.Errorf("failed to marshal public key: %v", pkiErr)
		return
	}
	ski := sha512.Sum512_256(sPki)

	tmpl := &x509.Certificate{
		SerialNumber: new(big.Int).SetBytes(serial[:20]),
		Subject: pkix.Name{
			CommonName: serverName + "-server",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(validYears, 0, 0),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  false,
		SubjectKeyId:          ski[:],
		AuthorityKeyId:        caCert.SubjectKeyId,
	}

	// 4. Sign the certificate with CA
	certDER, createCertErr := x509.CreateCertificate(rand.Reader, tmpl, caCert, &serverKey.PublicKey, caKey)
	if createCertErr != nil {
		err = fmt.Errorf("failed to sign server certificate: %v", createCertErr)
		return
	}
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	keyBytes, keyBytesErr := x509.MarshalECPrivateKey(serverKey)
	if keyBytesErr != nil {
		err = fmt.Errorf("failed to marshal server private key: %v", keyBytesErr)
		return
	}

	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	// 5. Generate tls-crypt-v2 server key with openvpn execution
	tlsCryptV2ServerKey, err = GenerateTlsCryptV2ServerKey()
	if err != nil {
		err = fmt.Errorf("failed to generate tls-crypt-v2 key: %v", err)
		return
	}

	return
}
