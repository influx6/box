package pkg

import (
	"bytes"
	"errors"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/influx6/devapp/static"
	"github.com/influx6/faux/filesystem"
	"github.com/influx6/mnet/certificates"
	"github.com/influx6/moz/gen"

	uuid "github.com/satori/go.uuid"
)

// errors ....
var (
	ErrSecretRequired = errors.New("Required to provide secret")
)

// vars ...
var (
	RSAKeyStrength                  = 4096
	CertificateAuthorityReqLifetime = time.Hour * 8760       // 1 year certificate requests
	CertificateAuthorityLifetime    = time.Hour * (5 * 8760) // 5 years CA
	certificateService              = certificates.CertificateAuthorityService{
		KeyStrength: RSAKeyStrength,
		LifeTime:    CertificateAuthorityLifetime,
		Serials: certificates.SerialService{
			Length: 128,
		},
	}
)

// BoxCtrl defines giving data related to a initialized box server
// instance.
type BoxCtrl struct {
	ID          string `toml:"id"`
	Secret      string `toml:"-"`
	Serial      string `toml:"serial"`
	Company     string `toml:"company"`
	ServerName  string `toml:"serverName"`
	BaseFiles   filesystem.FilePortal
	CaAuthority *certificates.CertificateAuthority
	CaRequests  *certificates.CertificateRequestService
	CaServer    *certificates.CertificateRequest
}

// NewBoxCtrl returns a new instance of BoxCtrl with appropriate values initialized.
func NewBoxCtrl(secret string, company string, serverName string, base filesystem.FilePortal) *BoxCtrl {
	if serverName == "" {
		serverName = "*"
	}

	return &BoxCtrl{
		BaseFiles:  base,
		Company:    company,
		Secret:     secret,
		ServerName: serverName,
	}
}

// LoadConfiguration attempts to load configuration data from received
// filesystem portal.
func (box *BoxCtrl) LoadConfiguration() error {
	if box.Secret == "" {
		return ErrSecretRequired
	}

	boxData, err := box.BaseFiles.Get("boxfile")
	if err != nil {
		return box.setupConfiguration()
	}

	if _, err := toml.Decode(string(boxData), box); err != nil {
		return err
	}

	return box.loadCertificateAuthority()
}

// setupConfiguration sets up necessary profiles and certificate files
// for giving box.
func (box *BoxCtrl) setupConfiguration() error {
	box.ID = uuid.NewV4().String()
	config := gen.SourceTextWithName(
		"box.createProfile",
		static.MustReadFile("configs/box.tml", true),
		nil,
		box,
	)

	var boxfile bytes.Buffer
	if _, err := config.WriteTo(&boxfile); err != nil {
		return err
	}

	if err := box.BaseFiles.Save("boxfile", boxfile.Bytes()); err != nil {
		return err
	}

	return box.loadCertificateAuthority()
}

// loadCertificateAuthority loads certificates for box else creating them
// which is used for issuing certificates to clients.
func (box *BoxCtrl) loadCertificateAuthority() error {
	certFiles, err := box.BaseFiles.Within("ca")
	if err != nil {
		return err
	}

	var ca certificates.CertificateAuthority
	if err := ca.Load(certFiles); err != nil {
		return box.setupCertificateAuthority()
	}

	return nil
}

// setupCertificateAuthority setups necessary root certificate for signing
// other ceritificates by box clients.
func (box *BoxCtrl) setupCertificateAuthority() error {
	certFiles, err := box.BaseFiles.Within("ca_authority")
	if err != nil {
		return err
	}

	// Remove all existing files.
	certFiles.RemoveAll()

	profile := certificates.CertificateProfile{
		Organization: box.Company,
		CommonName:   box.ServerName,
	}

	var requestService certificates.CertificateRequestService
	requestService.Profile = profile
	requestService.KeyStrength = RSAKeyStrength
	box.CaRequests = &requestService

	caService := certificateService
	caService.Profile = profile

	caAuthority, err := caService.New()
	if err != nil {
		return err
	}

	box.CaAuthority = &caAuthority
	if err := box.CaAuthority.Persist(certFiles); err != nil {
		return err
	}

	return box.loadServerCertificate()
}

// loadServerCertificate loads necessary certificates for setup server for
// box clients.
func (box *BoxCtrl) loadServerCertificate() error {
	certFiles, err := box.BaseFiles.Within("ca_server")
	if err != nil {
		return err
	}

	var serverCA certificates.CertificateRequest
	if err := serverCA.Load(certFiles); err != nil {
		return box.setupServerCertificate()
	}

	if serverCA.PrivateKey == nil {
		return errors.New("Invalid server certificate private key")
	}

	box.CaServer = &serverCA
	return nil
}

// setupServerCertificate sets up a new certificate which will be used for
// tls on a box server.
func (box *BoxCtrl) setupServerCertificate() error {
	certFiles, err := box.BaseFiles.Within("ca_server")
	if err != nil {
		return err
	}

	// RemoveAll certificate files.
	certFiles.RemoveAll()

	serverCA, err := box.CaRequests.New()
	if err != nil {
		return err
	}

	if err := box.CaAuthority.ApproveServerCertificateSigningRequest(
		&serverCA,
		certificateService.Serials,
		CertificateAuthorityReqLifetime,
	); err != nil {
		return err
	}

	box.CaServer = &serverCA
	return serverCA.Persist(certFiles)
}
