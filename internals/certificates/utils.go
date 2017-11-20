package certificates

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// decodePEMPrivateKey decodes byte slice has a pem file which contains
// a private key.
func decodePEMPrivateKey(keyData []byte) (*rsa.PrivateKey, error) {
	certkeyblock, _ := pem.Decode(keyData)
	if certkeyblock == nil {
		return nil, ErrInvalidPemBlock
	}

	if certkeyblock.Type != certKeyName {
		return nil, ErrInvalidCAKeyBlockType
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(certkeyblock.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// decodePEMCertificate decodes byte slice has a pem file which contains
// a signed x509.Certificate.
func decodePEMCertificate(certData []byte, header string) (*x509.Certificate, error) {
	certblock, _ := pem.Decode(certData)
	if certblock == nil {
		return nil, ErrInvalidPemBlock
	}

	if certblock.Type != header {
		return nil, ErrInvalidCABlockType
	}

	certificate, err := x509.ParseCertificate(certblock.Bytes)
	if err != nil {
		return nil, err
	}
	return certificate, nil
}
