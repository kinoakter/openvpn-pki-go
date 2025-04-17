package pki

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/google/uuid"
	"github.com/kinoakter/openvpn-pki-go/internal/db"
	"math/big"
	"time"
)

func CreateCA(name string) error {
	// Generate ECDSA key (secp521r1)
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate ECDSA key: %v", err)
	}

	// Create certificate template
	serial := sha512.Sum512([]byte(uuid.New().String())) // large serial based on UUID hash
	template := &x509.Certificate{
		SerialNumber: new(big.Int).SetBytes(serial[:20]),
		Subject:      pkix.Name{CommonName: name},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0), // valid for 10 years

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Self-sign the certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create CA certificate: %v", err)
	}

	// Encode cert and key as PEM
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %v", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	// Store in DB
	_, err = db.DB.Exec(context.Background(),
		`INSERT INTO ca (id, name, certificate, private_key) VALUES ($1, $2, $3, $4)`,
		uuid.New(), name, string(certPEM), string(keyPEM),
	)

	if err != nil {
		return fmt.Errorf("failed to store CA in database: %v", err)
	}

	return nil
}
