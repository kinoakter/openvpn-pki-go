// Package pki provides utilities for managing Public Key Infrastructure (PKI) operations,
// including certificate parsing, generation, and management for OpenVPN infrastructure.
package pki

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
)

// parseCertificate decodes and parses a PEM-encoded X.509 certificate.
//
// Parameters:
//   - certPEM: PEM-encoded certificate string
//
// Returns:
//   - cert: Parsed X.509 certificate structure
//   - err: Error if certificate parsing fails
func parseCertificate(certPEM string) (cert *x509.Certificate, err error) {
	block, _ := pem.Decode([]byte(certPEM))
	cert, err = x509.ParseCertificate(block.Bytes)

	return
}

// parseECPrivateKey decodes and parses a PEM-encoded ECDSA private key.
//
// Parameters:
//   - keyPEM: PEM-encoded private key string
//
// Returns:
//   - caKey: Parsed ECDSA private key structure
//   - err: Error if private key parsing fails
func parseECPrivateKey(keyPEM string) (caKey *ecdsa.PrivateKey, err error) {
	keyBlock, _ := pem.Decode([]byte(keyPEM))
	caKey, err = x509.ParseECPrivateKey(keyBlock.Bytes)

	return
}
