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
	"math/big"
	"time"
)

/*
var rootCert = x509.Certificate{
	Subject: pkix.Name{
		Country:      []string{"ID"},
		Organization: []string{"YK Brothers Co."},
		CommonName:   "Yeka Root CA",
	},
	NotBefore:   time.Now(),
	NotAfter:    time.Now().AddDate(10, 0, 0),
	IPAddresses: []net.IP{},
	DNSNames:    []string{},
}

var intermediateCert = x509.Certificate{
	Subject: pkix.Name{
		Country:      []string{"ID"},
		Organization: []string{"YK Intermediate Corp"},
		CommonName:   "Yeka Intermediate CA",
	},
	NotBefore:   time.Now(),
	NotAfter:    time.Now().AddDate(10, 0, 0),
	IPAddresses: []net.IP{},
	DNSNames:    []string{},
}

var serverCert = x509.Certificate{
	Subject: pkix.Name{
		Country:      []string{"ID"},
		Organization: []string{"Go Web Corp."},
		CommonName:   "Go Web",
	},
	NotBefore:   time.Now(),
	NotAfter:    time.Now().AddDate(10, 0, 0),
	IPAddresses: []net.IP{},
	DNSNames:    []string{"cc.local"},
}
*/

func CreateRootCA(cert *x509.Certificate) ([]byte, PrivateKey, error) {
	*cert = RootCert(*cert)
	return CreateCertificate(cert, nil, nil)
}

func CreateServerCertificate(cert, parent *x509.Certificate, parentKey PrivateKey) ([]byte, PrivateKey, error) {
	*cert = ServerCert(*cert)
	return CreateCertificate(cert, parent, parentKey)
}

// CreateCertificate creates a server certificate using ECDSA private key
// To create RootCA, parent & parentKey should be nill
func CreateCertificate(cert, parent *x509.Certificate, parentKey PrivateKey) ([]byte, PrivateKey, error) {
	key, err := GenerateECDSAKey()
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
	return rsa.GenerateKey(rand.Reader, 2048)
}

// ====================
// ====== Helper ======
// ====================

func RootCert(cert x509.Certificate) x509.Certificate {
	cert.NotBefore = time.Now()
	cert.SerialNumber = big.NewInt(1) // TODO: Randomize
	cert.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageCRLSign
	cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	cert.BasicConstraintsValid = true
	cert.IsCA = true
	cert.MaxPathLen = 2
	return cert
}

func IntermediateCert(cert x509.Certificate) x509.Certificate {
	cert.NotBefore = time.Now()
	cert.SerialNumber = big.NewInt(2) // TODO: Randomize
	cert.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageCRLSign
	cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	cert.BasicConstraintsValid = true
	cert.IsCA = true
	cert.MaxPathLen = 1
	return cert
}

func ServerCert(cert x509.Certificate) x509.Certificate {
	cert.NotBefore = time.Now()
	cert.SerialNumber = big.NewInt(3) // TODO: Randomize
	cert.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	cert.BasicConstraintsValid = true
	return cert
}

func WriteCert(cert []byte, w io.Writer) error {
	err := pem.Encode(w, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	return err
}

func WriteKey(key PrivateKey, w io.Writer) error {
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	err = pem.Encode(w, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	return err
}

func LoadCert(r io.Reader) (*x509.Certificate, error) {
	certBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	certPem, _ := pem.Decode(certBytes)
	return x509.ParseCertificate(certPem.Bytes)
}

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
