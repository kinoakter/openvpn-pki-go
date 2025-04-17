package pki

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/kinoakter/openvpn-pki-go/internal/db"
)

func IssueClientCertificate(commonName string) (certPEM, keyPEM, tlsCryptKey string, err error) {
	ctx := context.Background()

	// 1. Load CA cert and key from DB
	row := db.DB.QueryRow(ctx, `SELECT certificate, private_key FROM ca LIMIT 1`)
	var caCertPEM, caKeyPEM string
	if err := row.Scan(&caCertPEM, &caKeyPEM); err != nil {
		return "", "", "", fmt.Errorf("failed to load CA: %v", err)
	}

	// Decode CA certificate
	block, _ := pem.Decode([]byte(caCertPEM))
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse CA cert: %v", err)
	}

	// Decode CA private key
	keyBlock, _ := pem.Decode([]byte(caKeyPEM))
	caKey, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse CA key: %v", err)
	}

	// 2. Generate a new client keypair
	clientKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate client key: %v", err)
	}

	// 3. Create certificate template
	serial := sha512.Sum512([]byte(uuid.New().String()))
	tmpl := &x509.Certificate{
		SerialNumber: new(big.Int).SetBytes(serial[:20]),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0), // 1 year

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	// 4. Sign the certificate
	der, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, &clientKey.PublicKey, caKey)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to sign client cert: %v", err)
	}

	// 5. Encode to PEM
	certPEM = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
	privBytes, _ := x509.MarshalECPrivateKey(clientKey)
	keyPEM = string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}))

	// 6. Generate tls-crypt-v2 client key (just a secure random blob base64-encoded)
	rawKey := make([]byte, 256) // OpenVPN recommends 256 bytes
	if _, err := rand.Read(rawKey); err != nil {
		return "", "", "", fmt.Errorf("failed to generate tls-crypt-v2 key: %v", err)
	}
	tlsCryptKey = base64.StdEncoding.EncodeToString(rawKey)

	// 7. Store in DB
	_, err = db.DB.Exec(ctx,
		`INSERT INTO client_certificates (id, common_name, certificate, private_key, tls_crypt_v2_key) 
		 VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), commonName, certPEM, keyPEM, tlsCryptKey,
	)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to save client cert to DB: %v", err)
	}

	return certPEM, keyPEM, tlsCryptKey, nil
}
