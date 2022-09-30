package certs

import (
	"bytes"
	"crypto/x509"
	"errors"
	"time"
)

type Data struct {
	ID               int
	ParentID         int
	CertificateBytes []byte
	PrivateKeyBytes  []byte

	Info
}

type Info struct {
	CommonName   string    `json:"common_name"`
	Country      string    `json:"country"`
	Organization string    `json:"organization"`
	IPAddresses  []string  `json:"ip"`
	DNSNames     []string  `json:"dns"`
	ExpiredAt    time.Time `json:"expired_at"`
}

func (d Data) CertPair() (*x509.Certificate, PrivateKey, error) {
	if len(d.CertificateBytes) == 0 {
		return nil, nil, errors.New("empty certificate bytes")
	}

	rootCert, err := LoadCert(bytes.NewReader(d.CertificateBytes))
	if err != nil {
		return nil, nil, err
	}

	rootKey, err := LoadKey(bytes.NewReader(d.PrivateKeyBytes))
	return rootCert, rootKey, err
}

func (d *Data) SetCertPair(certBytes []byte, certKey PrivateKey) error {
	bufC := bytes.Buffer{}
	if err := WriteCert(certBytes, &bufC); err != nil {
		return err
	}

	bufK := bytes.Buffer{}
	if err := WriteKey(certKey, &bufK); err != nil {
		return err
	}

	d.CertificateBytes = bufC.Bytes()
	d.PrivateKeyBytes = bufK.Bytes()
	return nil
}
