package certificates

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

// const defines series of constant values
const (
	defaultSerialLength uint = 128
)

var (
	// ModernCiphers defines a list of modern tls cipher suites.
	ModernCiphers = []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	}
)

// SerialService defines a single method which always returns a serial number
// of length 128 by default.
type SerialService struct {
	Length uint
}

// New returns a new serial number acording to provided limit.
func (s SerialService) New() (*big.Int, error) {
	limit := defaultSerialLength
	if s.Length > 0 {
		limit = s.Length
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), limit)
	return rand.Int(rand.Reader, serialNumberLimit)
}

// PersistenceStore defines an interface which exposes a single method
// to persist a giving data into an underline store.
type PersistenceStore interface {
	Persist(string, []byte) error
	Retrieve(string) ([]byte, error)
}

// CertificateAuthority defines a struct which contains a generated certificate template with
// associated private and public keys.
type CertificateAuthority struct {
	PrivateKey  *rsa.PrivateKey
	PublicKey   *rsa.PublicKey
	Template    x509.Certificate
	Certificate *x509.Certificate
}

// Persist persist giving certificate into underline store.
func (ca CertificateAuthority) Persist(store PersistenceStore) error {
	var certBuffer bytes.Buffer
	if err := pem.Encode(&certBuffer, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca.Certificate.Raw,
	}); err != nil {
		return err
	}

	var keyBuffer bytes.Buffer
	if err := pem.Encode(&keyBuffer, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(ca.PrivateKey),
	}); err != nil {
		return err
	}

	if err := store.Persist("ca.crt", certBuffer.Bytes()); err != nil {
		return err
	}

	if err := store.Persist("ca.key", keyBuffer.Bytes()); err != nil {
		return err
	}

	return nil
}

// CertificateProfile holds authority profile data which are used to
// annotate a CA.
type CertificateProfile struct {
	Organization string `json:"org"`
	Country      string `json:"country"`
	Province     string `json:"province"`
	Local        string `json:"local"`
	Address      string `json:"address"`
	Postal       string `json:"postal"`
	SerialNumber string `json:"serial_number"`
	CommonName   string `json:"common_name"`
}

// CertificateAuthorityService generates a certificate with associated private key
// and public key, which can be saved into a giving persistent layer when supplied to
// it's persist method
type CertificateAuthorityService struct {
	KeyStrength int
	LifeTime    time.Duration
	Profile     CertificateProfile
	Serials     SerialService
	KeyUsages   []x509.ExtKeyUsage
	Emails      []string
	IPs         []string

	// General list of DNSNames for certificate.
	DNSNames []string

	// DNSNames to be excluded.
	ExDNSNames []string

	// DNSNames to be permitted.
	PermDNSNames []string
}

// New returns a new instance of CertificateAuthorty which implements the
// the necessary interface to write given certificate data into memory or
// into a given store.
func (cas CertificateAuthorityService) New() (CertificateAuthority, error) {
	var ca CertificateAuthority

	serial, err := cas.Serials.New()
	if err != nil {
		return ca, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, cas.KeyStrength)
	if err != nil {
		return ca, err
	}

	ca.PrivateKey = privateKey
	ca.PublicKey = &privateKey.PublicKey

	var ips []net.IP

	for _, ip := range cas.IPs {
		ips = append(ips, net.ParseIP(ip))
	}

	before := time.Now()

	var profile pkix.Name
	profile.Organization = []string{cas.Profile.Organization}
	profile.Country = []string{cas.Profile.Country}
	profile.Province = []string{cas.Profile.Province}
	profile.Locality = []string{cas.Profile.Local}
	profile.StreetAddress = []string{cas.Profile.Address}
	profile.PostalCode = []string{cas.Profile.Postal}

	var template x509.Certificate
	template.IsCA = true
	template.IPAddresses = ips
	template.Subject = profile
	template.NotBefore = before
	template.SerialNumber = serial
	template.DNSNames = cas.DNSNames
	template.EmailAddresses = cas.Emails
	template.BasicConstraintsValid = true
	template.NotAfter = before.Add(cas.LifeTime)
	template.ExcludedDNSDomains = cas.ExDNSNames
	template.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	template.ExtKeyUsage = append([]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}, cas.KeyUsages...)

	if len(cas.PermDNSNames) != 0 {
		template.PermittedDNSDomainsCritical = true
		template.PermittedDNSDomains = cas.PermDNSNames
	}

	ca.Template = template

	certData, err := x509.CreateCertificate(rand.Reader, &ca.Template, &ca.Template, ca.PublicKey, ca.PrivateKey)
	if err != nil {
		return ca, err
	}

	parsedCertificate, err := x509.ParseCertificate(certData)
	if err != nil {
		return ca, err
	}

	ca.Certificate = parsedCertificate

	return ca, nil
}
