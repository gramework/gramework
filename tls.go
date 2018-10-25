package gramework

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"
)

const localhost = "localhost"

func selfSignedCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if clientHello.ServerName != localhost {
		return nil, fmt.Errorf("self-signed certificate for %q not permitted", clientHello.ServerName)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %s", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},

		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 1, 0),

		DNSNames: []string{localhost},

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	return &tls.Certificate{
		Certificate: [][]byte{cert},
		PrivateKey:  priv,
	}, nil
}
