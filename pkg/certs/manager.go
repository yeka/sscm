package certs

import (
	"crypto/x509"
	"errors"
	"net"
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

func (man Manager) Create(cert *Data) (err error) {
	ips := make([]net.IP, len(cert.IPAddresses))
	for i, v := range cert.IPAddresses {
		ips[i] = net.ParseIP(v)
	}

	// assemble certificate information
	var c x509.Certificate
	if cert.ParentID == 0 {
		c = DefaultRootCA()
	} else {
		c = DefaultCert()
	}

	c.Subject.CommonName = cert.CommonName
	c.Subject.Country = []string{cert.Country}
	c.Subject.Organization = []string{cert.Organization}
	c.IPAddresses = ips
	c.DNSNames = cert.DNSNames

	// load parent certificate
	var rootCert *x509.Certificate
	var rootKey PrivateKey
	if cert.ParentID > 0 {
		rootData, err := man.Storage.Load(cert.ParentID)
		if err != nil {
			return err
		}

		rootCert, rootKey, err = rootData.CertPair()
		if err != nil {
			return err
		}
	}

	// generate certificate pair
	certByte, certKey, err := CreateCertificate(&c, rootCert, rootKey)
	if err != nil {
		return err
	}
	if err := cert.SetCertPair(certByte, certKey); err != nil {
		return err
	}

	// store certificate information
	return man.Storage.Store(cert)
}

func (inf Info) Validate() error {
	return errors.New("not yet implemented")
}
