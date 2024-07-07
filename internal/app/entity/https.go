package entity

import "bytes"

// HTTPSCredentials Contains certificate and private key credentials
type HTTPSCredentials struct {
	Cert       bytes.Buffer
	PrivateKey bytes.Buffer
}
