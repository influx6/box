package certificates

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"time"
)

// const defines series of constant values
const (
	defaultSerialLength   uint = 128
	certFileName               = "ca.cert"
	certKeyFileName            = "ca.key"
	reqcertFileName            = "req_ca.cert"
	reqcertKeyFileName         = "req_ca.key"
	reqcertRootCAFileName      = "req_root_ca.key"
	certTypeName               = "CERTIFICATE"
	rootCertTypeName           = "ROOT_CERTIFICATE"
	certKeyName                = "RSA PRIVATE KEY"
)

// errors ...
var (
	ErrExcludedDNSName          = errors.New("excluded DNSName")
	ErrNoCertificate            = errors.New("has no certificate")
	ErrNoRootCACertificate      = errors.New("has no root CA certificate")
	ErrNoCertificateRequest     = errors.New("has no certificate request")
	ErrNoPrivateKey             = errors.New("has no private key")
	ErrWrongSignatureAlgorithmn = errors.New("incorrect signature algorithmn received")
	ErrInvalidPemBlock          = errors.New("pem.Decode found no pem.Block data")
	ErrInvalidCABlockType       = errors.New("pem.Block has invalid block header for ca cert")
	ErrInvalidCAKeyBlockType    = errors.New("pem.Block has invalid block header for ca key")
	ErrEmptyCARawSlice          = errors.New("CA Raw slice is empty")
	ErrInvalidRawLength         = errors.New("CA Raw slice length is invalid")
	ErrInvalidRequestRawLength  = errors.New("RequestCA Raw slice length is invalid")
	ErrInvalidRootCARawLength   = errors.New("RootCA Raw slice length is invalid")
	ErrInvalidRawCertLength     = errors.New("Cert raw slice length is invalid")
	ErrInvalidRawCertKeyLength  = errors.New("Cert Key raw slice length is invalid")
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
// @mock
type PersistenceStore interface {
	Persist(string, []byte) error
	Retrieve(string) ([]byte, error)
}

// CertificateAuthority defines a struct which contains a generated certificate template with
// associated private and public keys.
type CertificateAuthority struct {
	PrivateKey  *rsa.PrivateKey
	PublicKey   *rsa.PublicKey
	Certificate *x509.Certificate
}

// SecondaryCertificateAuthority defines a certificate authority which is not a CA and is signed
// by a root CA.
type SecondaryCertificateAuthority struct {
	RootCA      *x509.Certificate
	Certificate *x509.Certificate
}

// RootCertificateRaw returns the raw version of the certificate.
func (sca SecondaryCertificateAuthority) RootCertificateRaw() ([]byte, error) {
	if sca.RootCA == nil {
		return nil, ErrNoRootCACertificate
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  rootCertTypeName,
		Bytes: sca.RootCA.Raw,
	}), nil
}

// CertificateRaw returns the raw version of the certificate.
func (sca SecondaryCertificateAuthority) CertificateRaw() ([]byte, error) {
	if sca.Certificate == nil {
		return nil, ErrNoCertificate
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  certTypeName,
		Bytes: sca.Certificate.Raw,
	}), nil
}

// Raw returns the whole CertificateAuthorty has a combined slice of bytes, with
// total length and length of parts added at the beginning.
// Format: [LENGTH OF ALL DATA][2] [LENGTH OF CERTIFICATE RAW][2] [LENGTH OF ROOTCA KEY][2] [CERTIFICIATE][ROOTCA]
// The above format then can be pulled and split properly to ensure matching data.
func (sca SecondaryCertificateAuthority) Raw() ([]byte, error) {
	var rootCABuffer bytes.Buffer
	if err := pem.Encode(&rootCABuffer, &pem.Block{
		Type:  rootCertTypeName,
		Bytes: sca.RootCA.Raw,
	}); err != nil {
		return nil, err
	}

	var certBuffer bytes.Buffer
	if err := pem.Encode(&certBuffer, &pem.Block{
		Type:  certKeyName,
		Bytes: sca.Certificate.Raw,
	}); err != nil {
		return nil, err
	}

	rootcertLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(rootcertLen, uint16(rootCABuffer.Len()))

	certLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(certLen, uint16(certBuffer.Len()))

	bitsLen := len(rootcertLen) + len(certLen)
	contentLen := certBuffer.Len() + rootCABuffer.Len()

	rawLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(rawLen, uint16(contentLen+bitsLen))

	var raw bytes.Buffer
	raw.Write(rawLen)
	raw.Write(certLen)
	raw.Write(rootcertLen)
	certBuffer.WriteTo(&raw)
	rootCABuffer.WriteTo(&raw)

	return raw.Bytes(), nil
}

// ApproveServerClientCertificateSigningRequest processes the provided CertificateRequest returning a new CertificateAuthorty
// which has being signed by this root CA.
// All received signed by this method receive ExtKeyUsageServerAuth and ExtKeyUsageClientAuth.
func (ca CertificateAuthority) ApproveServerClientCertificateSigningRequest(req *CertificateRequest, serial SerialService, lifeTime time.Duration) error {
	var secondaryCA SecondaryCertificateAuthority

	usage := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	template, err := ca.initCertificateRequest(req, serial, lifeTime, usage)
	if err != nil {
		return err
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, template, ca.Certificate, template.PublicKey, ca.PrivateKey)
	if err != nil {
		return err
	}

	certificate, err := x509.ParseCertificate(certificateBytes)
	if err != nil {
		return err
	}

	secondaryCA.Certificate = certificate
	secondaryCA.RootCA = ca.Certificate

	if err := req.ValidateAndAccept(secondaryCA, usage); err != nil {
		return err
	}

	return nil
}

// ApproveServerCertificateSigningRequest processes the provided CertificateRequest returning a new CertificateAuthorty
// which has being signed by this root CA.
// All received signed by this method receive ExtKeyUsageServerAuth alone.
func (ca CertificateAuthority) ApproveServerCertificateSigningRequest(req *CertificateRequest, serial SerialService, lifeTime time.Duration) error {
	var secondaryCA SecondaryCertificateAuthority

	usage := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	template, err := ca.initCertificateRequest(req, serial, lifeTime, usage)
	if err != nil {
		return err
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, template, ca.Certificate, template.PublicKey, ca.PrivateKey)
	if err != nil {
		return err
	}

	certificate, err := x509.ParseCertificate(certificateBytes)
	if err != nil {
		return err
	}

	secondaryCA.Certificate = certificate
	secondaryCA.RootCA = ca.Certificate

	if err := req.ValidateAndAccept(secondaryCA, usage); err != nil {
		return err
	}

	return nil
}

// ApproveClientCertificateSigningRequest processes the provided CertificateRequest returning a new CertificateAuthorty
// which has being signed by this root CA.
// All received signed by this method receive ExtKeyUsageClientAuth alone.
func (ca CertificateAuthority) ApproveClientCertificateSigningRequest(req *CertificateRequest, serial SerialService, lifeTime time.Duration) error {
	var secondaryCA SecondaryCertificateAuthority

	usage := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	template, err := ca.initCertificateRequest(req, serial, lifeTime, usage)
	if err != nil {
		return err
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, template, ca.Certificate, template.PublicKey, ca.PrivateKey)
	if err != nil {
		return err
	}

	certificate, err := x509.ParseCertificate(certificateBytes)
	if err != nil {
		return err
	}

	secondaryCA.RootCA = ca.Certificate
	secondaryCA.Certificate = certificate

	if err := req.ValidateAndAccept(secondaryCA, usage); err != nil {
		return err
	}

	return nil
}

func (ca CertificateAuthority) initCertificateRequest(creq *CertificateRequest, serial SerialService, lifeTime time.Duration, usages []x509.ExtKeyUsage) (*x509.Certificate, error) {
	serialNumber, err := serial.New()
	if err != nil {
		return nil, err
	}

	before := time.Now()
	req := creq.Request

	var template x509.Certificate
	template.SerialNumber = serialNumber
	template.Signature = req.Signature
	template.SignatureAlgorithm = req.SignatureAlgorithm
	template.PublicKey = req.PublicKey
	template.PublicKeyAlgorithm = req.PublicKeyAlgorithm
	template.Subject = req.Subject
	template.Issuer = ca.Certificate.Subject
	template.NotBefore = before
	template.NotAfter = before.Add(lifeTime)
	template.KeyUsage = x509.KeyUsageDigitalSignature
	template.DNSNames = req.DNSNames
	template.IPAddresses = req.IPAddresses
	template.EmailAddresses = req.EmailAddresses
	template.ExtKeyUsage = usages
	template.Extensions = req.Extensions
	template.ExtraExtensions = req.ExtraExtensions

	return &template, nil
}

// PrivateKeyRaw returns the raw version of the certificate's private key.
func (ca CertificateAuthority) PrivateKeyRaw() ([]byte, error) {
	if ca.PrivateKey == nil {
		return nil, ErrNoPrivateKey
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  certKeyName,
		Bytes: x509.MarshalPKCS1PrivateKey(ca.PrivateKey),
	}), nil
}

// CertificateRaw returns the raw version of the certificate.
func (ca CertificateAuthority) CertificateRaw() ([]byte, error) {
	if ca.Certificate == nil {
		return nil, ErrNoCertificate
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  certTypeName,
		Bytes: ca.Certificate.Raw,
	}), nil
}

// Raw returns the whole CertificateAuthorty has a combined slice of bytes, with
// total length and length of parts added at the beginning.
// Format: [LENGTH OF ALL DATA][2] [LENGTH OF CERTIFICATE RAW][2] [LENGTH OF PRIVATE KEY][2] [CERTIFICIATE][PRIVIATEKEY]
// The above format then can be pulled and split properly to ensure matching data.
func (ca CertificateAuthority) Raw() ([]byte, error) {
	var certBuffer bytes.Buffer
	if err := pem.Encode(&certBuffer, &pem.Block{
		Type:  certTypeName,
		Bytes: ca.Certificate.Raw,
	}); err != nil {
		return nil, err
	}

	var keyBuffer bytes.Buffer
	if err := pem.Encode(&keyBuffer, &pem.Block{
		Type:  certKeyName,
		Bytes: x509.MarshalPKCS1PrivateKey(ca.PrivateKey),
	}); err != nil {
		return nil, err
	}

	certLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(certLen, uint16(certBuffer.Len()))

	keyLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(keyLen, uint16(keyBuffer.Len()))

	bitsLen := len(keyLen) + len(certLen)
	contentLen := certBuffer.Len() + keyBuffer.Len()

	rawLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(rawLen, uint16(contentLen+bitsLen))

	var raw bytes.Buffer
	raw.Write(rawLen)
	raw.Write(certLen)
	raw.Write(keyLen)
	certBuffer.WriteTo(&raw)
	keyBuffer.WriteTo(&raw)

	return raw.Bytes(), nil
}

// FromRaw decodes the whole raw data into CertificateAuthorty, using the format below
// Format: [LENGTH OF ALL DATA][2] [LENGTH OF CERTIFICATE RAW][2] [LENGTH OF PRIVATE KEY][2] [CERTIFICIATE][PRIVIATEKEY]
// The above format then can be pulled and split properly to ensure matching data.
func (ca *CertificateAuthority) FromRaw(raw []byte) error {
	if len(raw) == 0 {
		return ErrEmptyCARawSlice
	}

	if len(raw) <= 7 {
		return ErrInvalidRawLength
	}

	rawLenBytes := raw[0:2]
	certLenBytes := raw[2:4]
	keyLenBytes := raw[4:6]
	rest := raw[6:]

	rawLen := int(binary.LittleEndian.Uint16(rawLenBytes))
	if len(raw[2:]) != rawLen {
		return ErrInvalidRawLength
	}

	keyLen := int(binary.LittleEndian.Uint16(keyLenBytes))
	certLen := int(binary.LittleEndian.Uint16(certLenBytes))

	certKeyInfoLen := len(certLenBytes) + len(keyLenBytes)
	realRawLen := rawLen - certKeyInfoLen

	if (realRawLen - keyLen) != certLen {
		return ErrInvalidRawCertLength
	}

	if (realRawLen - certLen) != keyLen {
		return ErrInvalidRawCertKeyLength
	}

	certRaw := rest[:certLen]
	if len(certRaw) != certLen {
		return ErrInvalidRawCertLength
	}

	keyRaw := rest[certLen:]
	if len(keyRaw) != keyLen {
		return ErrInvalidRawCertKeyLength
	}

	certblock, _ := pem.Decode(certRaw)
	if certblock == nil {
		return ErrInvalidPemBlock
	}

	if certblock.Type != certTypeName {
		return ErrInvalidCABlockType
	}

	certificate, err := x509.ParseCertificate(certblock.Bytes)
	if err != nil {
		return err
	}

	ca.Certificate = certificate

	certkeyblock, _ := pem.Decode(keyRaw)
	if certkeyblock == nil {
		return ErrInvalidPemBlock
	}

	if certkeyblock.Type != certKeyName {
		return ErrInvalidCAKeyBlockType
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(certkeyblock.Bytes)
	if err != nil {
		return err
	}

	ca.PrivateKey = privateKey
	ca.PublicKey = &privateKey.PublicKey

	return nil
}

// Load loads certificate and key from provided PersistenceStore.
func (ca *CertificateAuthority) Load(store PersistenceStore) error {
	cert, err := store.Retrieve(certFileName)
	if err != nil {
		return err
	}

	certificate, err := decodePEMCertificate(cert, certTypeName)
	if err != nil {
		return err
	}

	ca.Certificate = certificate

	certKey, err := store.Retrieve(certKeyFileName)
	if err != nil {
		return err
	}

	privateKey, err := decodePEMPrivateKey(certKey)
	if err != nil {
		return err
	}

	ca.PrivateKey = privateKey
	ca.PublicKey = &privateKey.PublicKey
	return nil
}

// Persist persist giving certificate into underline store.
func (ca CertificateAuthority) Persist(store PersistenceStore) error {
	certBytes, err := ca.CertificateRaw()
	if err != nil {
		return err
	}

	keyBytes, err := ca.PrivateKeyRaw()
	if err != nil {
		return err
	}

	if err := store.Persist(certFileName, certBytes); err != nil {
		return err
	}

	if err := store.Persist(certKeyFileName, keyBytes); err != nil {
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
	CommonName   string `json:"common_name"`
}

// CertificateAuthorityService generates a certificate with associated private key
// and public key, which can be saved into a giving persistent layer when supplied to
// it's persist method
type CertificateAuthorityService struct {
	Version     int
	KeyStrength int
	LifeTime    time.Duration
	// SignatureAlgorithm x509.SignatureAlgorithm
	Profile   CertificateProfile
	Serials   SerialService
	KeyUsages []x509.ExtKeyUsage
	Emails    []string
	IPs       []string

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
	profile.CommonName = cas.Profile.CommonName
	profile.Organization = []string{cas.Profile.Organization}
	profile.Country = []string{cas.Profile.Country}
	profile.Province = []string{cas.Profile.Province}
	profile.Locality = []string{cas.Profile.Local}
	profile.StreetAddress = []string{cas.Profile.Address}
	profile.PostalCode = []string{cas.Profile.Postal}

	var template x509.Certificate
	template.Version = cas.Version
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
	template.SignatureAlgorithm = x509.SHA256WithRSA
	// template.SignatureAlgorithm = cas.SignatureAlgorithm
	template.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	template.ExtKeyUsage = append([]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}, cas.KeyUsages...)

	if len(cas.PermDNSNames) != 0 {
		template.PermittedDNSDomainsCritical = true
		template.PermittedDNSDomains = cas.PermDNSNames
	}

	certData, err := x509.CreateCertificate(rand.Reader, &template, &template, ca.PublicKey, ca.PrivateKey)
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

// CertificateRequest defines a struct which contains a generated certificate request template with
// associated private and public keys.
type CertificateRequest struct {
	PrivateKey  *rsa.PrivateKey
	PublicKey   *rsa.PublicKey
	Request     *x509.CertificateRequest
	SecondaryCA SecondaryCertificateAuthority
}

// RequestRaw returns the raw bytes that make up the request.
func (ca CertificateRequest) RequestRaw() ([]byte, error) {
	if ca.Request == nil {
		return nil, ErrNoCertificateRequest
	}
	return ca.Request.Raw, nil
}

// PrivateKeyRaw returns the raw version of the certificate's private key.
func (ca CertificateRequest) PrivateKeyRaw() ([]byte, error) {
	if ca.PrivateKey == nil {
		return nil, ErrNoPrivateKey
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  certKeyName,
		Bytes: x509.MarshalPKCS1PrivateKey(ca.PrivateKey),
	}), nil
}

// Load loads certificate and key from provided PersistenceStore.
func (ca *CertificateRequest) Load(store PersistenceStore) error {
	cert, err := store.Retrieve(reqcertFileName)
	if err != nil {
		return err
	}

	certificate, err := decodePEMCertificate(cert, certTypeName)
	if err != nil {
		return err
	}

	ca.SecondaryCA.Certificate = certificate

	certKey, err := store.Retrieve(reqcertKeyFileName)
	if err != nil {
		return err
	}

	privateKey, err := decodePEMPrivateKey(certKey)
	if err != nil {
		return err
	}

	ca.PrivateKey = privateKey
	ca.PublicKey = &privateKey.PublicKey

	rootCACert, err := store.Retrieve(reqcertRootCAFileName)
	if err != nil {
		return err
	}

	rootCertificate, err := decodePEMCertificate(rootCACert, rootCertTypeName)
	if err != nil {
		return err
	}

	ca.SecondaryCA.RootCA = rootCertificate

	return nil
}

// Persist persist giving certificate into underline store.
func (ca CertificateRequest) Persist(store PersistenceStore) error {
	certBytes, err := ca.SecondaryCA.CertificateRaw()
	if err != nil {
		return err
	}

	rootCABytes, err := ca.SecondaryCA.RootCertificateRaw()
	if err != nil {
		return err
	}

	keyBytes, err := ca.PrivateKeyRaw()
	if err != nil {
		return err
	}

	if err := store.Persist(reqcertFileName, certBytes); err != nil {
		return err
	}

	if err := store.Persist(reqcertRootCAFileName, rootCABytes); err != nil {
		return err
	}

	if err := store.Persist(reqcertKeyFileName, keyBytes); err != nil {
		return err
	}

	return nil
}

// ValidateAndAccept takes the provided request response and rootCA, validating the fact that the certifcate comes from the rootCA
// before setting the certificate has the certificate and setting the rootCA has it's RootCA. You must take care to ensure
// this incoming ones match the Certificate request data.
// It uses Sha256
func (ca *CertificateRequest) ValidateAndAccept(sec SecondaryCertificateAuthority, keyUsage []x509.ExtKeyUsage) error {
	if sec.Certificate.SignatureAlgorithm != x509.SHA256WithRSA {
		return ErrWrongSignatureAlgorithmn
	}

	// sha := sha256.New()
	// sha.Write(sec.Certificate.RawTBSCertificate)
	// hashed := sha.Sum(nil)
	//
	// if err := rsa.VerifyPKCS1v15(ca.PublicKey, crypto.SHA256, hashed[:], sec.Certificate.Signature); err != nil {
	// 	fmt.Printf("Verification Issue: %+q\n", err)
	// 	return err
	// }

	certpool := x509.NewCertPool()
	certpool.AddCert(sec.RootCA)

	options := x509.VerifyOptions{Roots: certpool, KeyUsages: keyUsage}
	if _, err := sec.Certificate.Verify(options); err != nil {
		return err
	}

	ca.SecondaryCA = sec
	return nil
}

// Raw returns the whole CertificateAuthorty has a combined slice of bytes, with
// total length and length of parts added at the beginning.
// Format:
// [LENGTH OF ALL DATA][2] [LENGTH OF CERTIFICATE REQUEST][2] [LENGTH OF CERTIFICATE RAW][2] [LENGTH OF ROOTCA KEY][2] [CERTIFICIATE][ROOTCA]
// The above format then can be pulled and split properly to ensure matching data.
func (ca CertificateRequest) Raw() ([]byte, error) {
	rootCA, err := ca.SecondaryCA.RootCertificateRaw()
	if err != nil {
		return nil, err
	}

	rootcertLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(rootcertLen, uint16(len(rootCA)))

	car, err := ca.SecondaryCA.CertificateRaw()
	if err != nil {
		return nil, err
	}

	certLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(certLen, uint16(len(car)))

	rq, err := ca.RequestRaw()
	if err != nil {
		return nil, err
	}

	reqLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(reqLen, uint16(len(rq)))

	bitsLen := len(rootcertLen) + len(certLen) + len(reqLen)
	contentLen := len(car) + len(rootCA) + len(rq)

	rawLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(rawLen, uint16(contentLen+bitsLen))

	var raw bytes.Buffer
	raw.Write(rawLen)
	raw.Write(reqLen)
	raw.Write(certLen)
	raw.Write(rootcertLen)
	raw.Write(rq)
	raw.Write(car)
	raw.Write(rootCA)

	return raw.Bytes(), nil
}

// FromRaw decodes the whole raw data into CertificateAuthorty, using the format below
// Format: [LENGTH OF ALL DATA][2] [LENGTH OF CERTIFICATE RAW][2] [LENGTH OF PRIVATE KEY][2] [CERTIFICIATE][PRIVIATEKEY]
// The above format then can be pulled and split properly to ensure matching data.
func (ca *CertificateRequest) FromRaw(raw []byte) error {
	if len(raw) == 0 {
		return ErrEmptyCARawSlice
	}

	if len(raw) <= 8 {
		return ErrInvalidRawLength
	}

	rawLenBytes := raw[0:2]
	reqLenBytes := raw[2:4]
	certLenBytes := raw[4:6]
	rootCertLenBytes := raw[6:8]
	rawData := raw[8:]

	rawLen := int(binary.LittleEndian.Uint16(rawLenBytes))
	if len(raw[2:]) != rawLen {
		return ErrInvalidRawLength
	}

	reqLen := int(binary.LittleEndian.Uint16(reqLenBytes))
	certLen := int(binary.LittleEndian.Uint16(certLenBytes))
	rootCertLen := int(binary.LittleEndian.Uint16(rootCertLenBytes))

	certKeyInfoLen := len(certLenBytes) + len(reqLenBytes) + len(rootCertLenBytes)
	realRawLen := rawLen - certKeyInfoLen

	if (realRawLen - reqLen) != (rootCertLen + certLen) {
		return ErrInvalidRawLength
	}

	req := rawData[:reqLen]
	if len(req) != reqLen {
		return ErrInvalidRequestRawLength
	}

	certLenTotal := certLen + reqLen
	cert := rawData[reqLen:certLenTotal]
	if len(cert) != certLen {
		return ErrInvalidRawCertLength
	}

	// rootCertLenTotal := certLen + rootCertLen
	rootCert := rawData[certLenTotal:]
	if len(rootCert) != rootCertLen {
		return ErrInvalidRootCARawLength
	}

	certificateRequest, err := x509.ParseCertificateRequest(req)
	if err != nil {
		return err
	}

	ca.Request = certificateRequest

	certificate, err := decodePEMCertificate(cert, certTypeName)
	if err != nil {
		return err
	}

	ca.SecondaryCA.Certificate = certificate

	rootCA, err := decodePEMCertificate(rootCert, rootCertTypeName)
	if err != nil {
		return err
	}

	ca.SecondaryCA.RootCA = rootCA
	return nil
}

// CertificateRequestService generates a certificate request with associated private key
// and public key, which can be sent over the wire or directly to a CeritificateAuthority
// for signing.
type CertificateRequestService struct {
	Version     int
	KeyStrength int
	Profile     CertificateProfile
	Emails      []string
	IPs         []string
	// SignatureAlgorithmn x509.SignatureAlgorithm

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
func (cas CertificateRequestService) New() (CertificateRequest, error) {
	var ca CertificateRequest

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

	var profile pkix.Name
	profile.CommonName = cas.Profile.CommonName
	profile.Organization = []string{cas.Profile.Organization}
	profile.Country = []string{cas.Profile.Country}
	profile.Province = []string{cas.Profile.Province}
	profile.Locality = []string{cas.Profile.Local}
	profile.StreetAddress = []string{cas.Profile.Address}
	profile.PostalCode = []string{cas.Profile.Postal}

	var template x509.CertificateRequest
	template.Version = cas.Version
	template.IPAddresses = ips
	template.Subject = profile
	template.DNSNames = cas.DNSNames
	template.EmailAddresses = cas.Emails
	template.SignatureAlgorithm = x509.SHA256WithRSA

	certData, err := x509.CreateCertificateRequest(rand.Reader, &template, ca.PrivateKey)
	if err != nil {
		return ca, err
	}

	parsedRequest, err := x509.ParseCertificateRequest(certData)
	if err != nil {
		return ca, err
	}

	ca.Request = parsedRequest

	return ca, nil
}
