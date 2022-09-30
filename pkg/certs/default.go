package certs

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func DefaultRootCA() x509.Certificate {
	return x509.Certificate{
		SerialNumber: big.NewInt(2019), // TODO: Randomize this
		Subject: pkix.Name{
			CommonName:    "Common Name",
			Organization:  []string{"Organization, INC."},
			Country:       []string{"Country"},
			Province:      []string{"Province"},
			Locality:      []string{"Locality"},
			StreetAddress: []string{"StreetAddress"},
			PostalCode:    []string{"Postal"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		// MaxPathLen:            2,
	}
}

func DefaultIntermediateCert() x509.Certificate {
	return x509.Certificate{
		SerialNumber: big.NewInt(1658), // TODO: Randomize this
		Subject: pkix.Name{
			CommonName:    "Common Name",
			Organization:  []string{"Organization, INC."},
			Country:       []string{"Country"},
			Province:      []string{"Province"},
			Locality:      []string{"Locality"},
			StreetAddress: []string{"StreetAddress"},
			PostalCode:    []string{"Postal"},
		},
		// IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// DNSNames:     []string{"ocbc.my", "raih.my", "nyala.my"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		// MaxPathLen:            1,
	}
}

func DefaultCert() x509.Certificate {
	return x509.Certificate{
		SerialNumber: big.NewInt(1658), // TODO: Randomize this
		Subject: pkix.Name{
			CommonName:    "Common Name",
			Organization:  []string{"Organization, INC."},
			Country:       []string{"Country"},
			Province:      []string{"Province"},
			Locality:      []string{"Locality"},
			StreetAddress: []string{"StreetAddress"},
			PostalCode:    []string{"Postal"},
		},
		// IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// DNSNames:     []string{"ocbc.my", "raih.my", "nyala.my"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature, //x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
		// BasicConstraintsValid: true,
	}
}
