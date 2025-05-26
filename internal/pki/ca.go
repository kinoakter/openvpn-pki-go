// Package pki provides Public Key Infrastructure (PKI) utilities for creating and managing
// certificates, including Certificate Authority (CA) certificates, server certificates,
// and client certificates for secure communication.
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
	"github.com/google/uuid"
	"math/big"
	"time"
)

const DefaultCAValidityYears = 3

// CreateCACert creates a new Certificate Authority (CA) certificate and its corresponding private key.
// It generates an ECDSA key pair using the P-521 curve and creates a self-signed CA certificate
// with the specified common name and validity period.
//
// Parameters:
//   - commonName: The CA certificate's subject common name
//   - validYears: The number of years the certificate will be valid (default to 3 if 0 is provided)
//
// Returns:
//   - certPEM: The CA certificate in PEM format
//   - keyPEM: The private key in PEM format
//   - err: Error if certificate creation fails
func CreateCACert(commonName string, validYears int) (certPEM, keyPEM []byte, err error) {
	if validYears == 0 {
		validYears = DefaultCAValidityYears // By default, 3 years.
	}

	// Generate ECDSA key (secp521r1)
	var privateKey *ecdsa.PrivateKey
	privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate ECDSA key: %v", err)
	}

	// Create certificate template
	serial := sha512.Sum512([]byte(uuid.New().String())) // large serial based on UUID hash
	var subjectKeyId = sha512.Sum512_256(serial[:])
	template := &x509.Certificate{
		SerialNumber:          new(big.Int).SetBytes(serial[:20]),
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now().UTC(),
		NotAfter:              time.Now().AddDate(validYears, 0, 0).UTC(),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		SubjectKeyId:          subjectKeyId[:],
	}

	// Self-sign the certificate
	var derBytes []byte
	derBytes, err = x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CA certificate: %v", err)
	}

	// Encode cert and key as PEM
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	var keyBytes []byte
	keyBytes, err = x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %v", err)
	}
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	return
}
