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

type Storage interface {
	// Store should override cert.ID using database ID on success
	Store(cert *Data) (err error)

	Load(id int) (data Data, err error)

	// Search certificates based on query string
	// empty query means just list the certificate
	// parentId = -1, list all certificates
	// parentId = 0, list root certificates
	// parentId > 0, list certificates under given parentId
	Search(query string, parentId int) ([]Data, error)
}

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

func (man Manager) AddRootCA(cert *Data) (err error) {
	ips := make([]net.IP, len(cert.IPAddresses))
	for i, v := range cert.IPAddresses {
		ips[i] = net.ParseIP(v)
	}
	certs := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cert.CommonName,
			Country:      []string{cert.Country},
			Organization: []string{cert.Organization},
		},
		NotBefore:   time.Now(),
		NotAfter:    cert.ExpiredAt,
		IPAddresses: ips,
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

	cert.ParentID = 0 // root
	cert.CertificateBytes = cb
	cert.PrivateKeyBytes = kb.Bytes()

	return man.Store(cert)
}

func (man Manager) AddServerCertificate(cert *Data, parentId int) (err error) {
	parent, err := man.Load(parentId)
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

	ips := make([]net.IP, len(cert.IPAddresses))
	for i, v := range cert.IPAddresses {
		ips[i] = net.ParseIP(v)
	}
	certs := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cert.CommonName,
			Country:      []string{cert.Country},
			Organization: []string{cert.Organization},
		},
		NotBefore:   time.Now(),
		NotAfter:    cert.ExpiredAt,
		IPAddresses: ips,
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

	cert.ParentID = parentId
	cert.CertificateBytes = b
	cert.PrivateKeyBytes = kb.Bytes()
	return man.Store(cert)
}

func (inf Info) Validate() error {
	return errors.New("not yet implemented")
}
