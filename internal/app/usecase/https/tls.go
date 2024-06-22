package https

import (
	"crypto/tls"
	"fmt"
)

// GenerateTLSCert Generates TLS sertificate using HTTPS credentials
func GenerateTLSCert() (tls.Certificate, error) {
	httpsCreds, err := GenerateHTTPSCredentials()
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("error while generating certificate: %w", err)
	}

	cert, err := tls.X509KeyPair(httpsCreds.Cert.Bytes(), httpsCreds.PrivateKey.Bytes())
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("error while generating TLS keys: %w", err)
	}

	return cert, nil
}
