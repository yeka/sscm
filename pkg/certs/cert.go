package certs

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"net"
	"time"
)

type Manager struct {
	Storage
}

type SearchMode int

const RootOnly SearchMode = 0
const NonRootOnly SearchMode = 1
const AllCertificates SearchMode = 2

type Storage interface {
	// Store should override cert.ID using database ID on success
	Store(cert *Data) (err error)
	Load(id int, isRoot bool) (data Data, err error)
	Search(query string, mode SearchMode) ([]Data, error)
}

type Data struct {
	ID               int
	IsRoot           bool
	CertificateBytes []byte
	PrivateKeyBytes  []byte

	Info
}

type Info struct {
	CommonName   string
	Country      string
	Organization string
	IPAddresses  []net.IP
	DNSNames     []string
	ExpiredAt    time.Time
}

func (man Manager) AddRootCA(cert *Data) (err error) {
	certs := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cert.CommonName,
			Country:      []string{cert.Country},
			Organization: []string{cert.Organization},
		},
		NotBefore:   time.Now(),
		NotAfter:    cert.ExpiredAt,
		IPAddresses: cert.IPAddresses,
		DNSNames:    cert.DNSNames,
	}

	cb, k, err := CreateRootCA(&certs)
	if err != nil {
		return err
	}
	kb := &bytes.Buffer{}
	err = WriteKey(k, kb)
	if err != nil {
		return err
	}

	cert.IsRoot = true
	cert.CertificateBytes = cb
	cert.PrivateKeyBytes = kb.Bytes()

	return man.Store(cert)
}

func (man Manager) AddServerCertificate(cert *Data, parentId int) (err error) {
	parent, err := man.Load(parentId, true)
	if err != nil {
		return err
	}
	pcerts, err := LoadCert(bytes.NewReader(parent.CertificateBytes))
	if err != nil {
		return err
	}
	pkey, err := LoadKey(bytes.NewReader(parent.PrivateKeyBytes))
	if err != nil {
		return err
	}

	certs := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cert.CommonName,
			Country:      []string{cert.Country},
			Organization: []string{cert.Organization},
		},
		NotBefore:   time.Now(),
		NotAfter:    cert.ExpiredAt,
		IPAddresses: cert.IPAddresses,
		DNSNames:    cert.DNSNames,
	}

	b, k, err := CreateServerCertificate(&certs, pcerts, pkey)
	if err != nil {
		return err
	}
	kb := &bytes.Buffer{}
	err = WriteKey(k, kb)
	if err != nil {
		return err
	}

	cert.IsRoot = false
	cert.CertificateBytes = b
	cert.PrivateKeyBytes = kb.Bytes()
	return man.Store(cert)
}

func (inf Info) Validate() error {
	return errors.New("not yet implemented")
}
