package certs

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
)

// CreateCertificate creates a server certificate using ECDSA private key
// To create RootCA, parent & parentKey should be nill
func CreateCertificate(cert, parent *x509.Certificate, parentKey PrivateKey) ([]byte, PrivateKey, error) {
	key, err := GenerateECDSAKey()
	// key, err := GenerateRSAKey()
	if err != nil {
		return nil, nil, err
	}

	if parent == nil {
		parent = cert
	}
	if parentKey == nil {
		parentKey = key
	}

	b, err := x509.CreateCertificate(rand.Reader, cert, parent, key.Public(), parentKey)
	return b, key, err
}

// ==================================
// ====== Interface & Function ======
// ==================================

type PrivateKey interface {
	Public() crypto.PublicKey
}

type PrivateKeyGenerator func() (PrivateKey, error)

func GenerateECDSAKey() (PrivateKey, error) {
	// ECDSA 256 is believed to be better than RSA (May 2022)
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func GenerateRSAKey() (PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

// ====================
// ====== Helper ======
// ====================

// WriteCert writes certificate bytes to an io.Writer
func WriteCert(cert []byte, w io.Writer) error {
	err := pem.Encode(w, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	return err
}

// WriteKey writes certificate key to an io.Writer
func WriteKey(key PrivateKey, w io.Writer) error {
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	err = pem.Encode(w, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	return err
}

// LoadCert loads a certificate from an io.Reader
func LoadCert(r io.Reader) (*x509.Certificate, error) {
	certBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	certPem, _ := pem.Decode(certBytes)
	return x509.ParseCertificate(certPem.Bytes)
}

// LoadKey loads a certificate key from an io.Reader
func LoadKey(r io.Reader) (PrivateKey, error) {
	keyBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	keyPem, _ := pem.Decode(keyBytes)
	key, err := x509.ParsePKCS8PrivateKey(keyPem.Bytes)
	if err != nil {
		return nil, err
	}
	return key.(PrivateKey), nil
}
