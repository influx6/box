package certificates_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/influx6/box/internals/certificates"
	"github.com/influx6/box/internals/certificates/mocks"
	"github.com/influx6/faux/tests"
)

func TestCertificateRequestService(t *testing.T) {
	storeMap := make(map[string][]byte)
	var store mocks.PersistenceStoreMock
	store.GetFunc, store.PersistFunc = mocks.MapStore(storeMap)

	serials := certificates.SerialService{Length: 128}
	profile := certificates.CertificateProfile{
		Local:        "Lagos",
		Organization: "DreamBench",
		CommonName:   "DreamBench Inc",
		Country:      "Nigeria",
		Province:     "South-West",
	}

	var service certificates.CertificateAuthorityService
	service.KeyStrength = 4096
	service.LifeTime = (time.Hour * 8760)
	service.Profile = profile
	service.Serials = serials
	service.Emails = append([]string{}, "alex.ewetumo@dreambench.io")

	ca, err := service.New()
	if err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateAuthority")
	}
	tests.Passed("Should have generated new CertificateAuthority")

	var requestService certificates.CertificateRequestService
	requestService.Profile = profile
	requestService.KeyStrength = 2048

	reqCA, err := requestService.New()
	if err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateRequest")
	}
	tests.Passed("Should have generated new CertificateRequest")

	// Generate Client and Server Auth Certificate.
	if err := ca.ApproveServerClientCertificateSigningRequest(&reqCA, serials, time.Hour*8760); err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateRequest")
	}
	tests.Passed("Should have generated new CertificateRequest")

	if reqCA.SecondaryCA.RootCA == nil {
		tests.FailedWithError(err, "Should have generated new Certificate for request")
	}
	tests.Passed("Should have generated new Certificate for request")
}

func TestCertificateRequestServiceForClient(t *testing.T) {
	storeMap := make(map[string][]byte)
	var store mocks.PersistenceStoreMock
	store.GetFunc, store.PersistFunc = mocks.MapStore(storeMap)

	serials := certificates.SerialService{Length: 128}
	profile := certificates.CertificateProfile{
		Local:        "Lagos",
		Organization: "DreamBench",
		CommonName:   "DreamBench Inc",
		Country:      "Nigeria",
		Province:     "South-West",
	}

	var service certificates.CertificateAuthorityService
	service.KeyStrength = 4096
	service.LifeTime = (time.Hour * 8760)
	service.Profile = profile
	service.Serials = serials
	service.Emails = append([]string{}, "alex.ewetumo@dreambench.io")

	ca, err := service.New()
	if err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateAuthority")
	}
	tests.Passed("Should have generated new CertificateAuthority")

	var requestService certificates.CertificateRequestService
	requestService.Profile = profile
	requestService.KeyStrength = 2048

	reqCA, err := requestService.New()
	if err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateRequest")
	}
	tests.Passed("Should have generated new CertificateRequest")

	// Generate Client and Server Auth Certificate.
	if err := ca.ApproveClientCertificateSigningRequest(&reqCA, serials, time.Hour*8760); err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateRequest")
	}
	tests.Passed("Should have generated new CertificateRequest")
}

func TestCertificateRequestServiceForServer(t *testing.T) {
	storeMap := make(map[string][]byte)
	var store mocks.PersistenceStoreMock
	store.GetFunc, store.PersistFunc = mocks.MapStore(storeMap)

	serials := certificates.SerialService{Length: 128}
	profile := certificates.CertificateProfile{
		Local:        "Lagos",
		Organization: "DreamBench",
		CommonName:   "DreamBench Inc",
		Country:      "Nigeria",
		Province:     "South-West",
	}

	var service certificates.CertificateAuthorityService
	service.KeyStrength = 4096
	service.LifeTime = (time.Hour * 8760)
	service.Profile = profile
	service.Serials = serials
	service.Emails = append([]string{}, "alex.ewetumo@dreambench.io")

	ca, err := service.New()
	if err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateAuthority")
	}
	tests.Passed("Should have generated new CertificateAuthority")

	var requestService certificates.CertificateRequestService
	requestService.Profile = profile
	requestService.KeyStrength = 2048

	reqCA, err := requestService.New()
	if err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateRequest")
	}
	tests.Passed("Should have generated new CertificateRequest")

	// Generate Client and Server Auth Certificate.
	if err := ca.ApproveServerCertificateSigningRequest(&reqCA, serials, time.Hour*8760); err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateRequest")
	}
	tests.Passed("Should have generated new CertificateRequest")
}

func TestCertificateRequestRawLoading(t *testing.T) {
	storeMap := make(map[string][]byte)
	var store mocks.PersistenceStoreMock
	store.GetFunc, store.PersistFunc = mocks.MapStore(storeMap)

	serials := certificates.SerialService{Length: 128}
	profile := certificates.CertificateProfile{
		Local:        "Lagos",
		Organization: "DreamBench",
		CommonName:   "DreamBench Inc",
		Country:      "Nigeria",
		Province:     "South-West",
	}

	var service certificates.CertificateAuthorityService
	service.KeyStrength = 4096
	service.LifeTime = (time.Hour * 8760)
	service.Profile = profile
	service.Serials = serials
	service.Emails = append([]string{}, "alex.ewetumo@dreambench.io")

	ca, err := service.New()
	if err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateAuthority")
	}
	tests.Passed("Should have generated new CertificateAuthority")

	var requestService certificates.CertificateRequestService
	requestService.Profile = profile
	requestService.KeyStrength = 2048

	reqCA, err := requestService.New()
	if err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateRequest")
	}
	tests.Passed("Should have generated new CertificateRequest")

	// Generate Client and Server Auth Certificate.
	if err := ca.ApproveServerClientCertificateSigningRequest(&reqCA, serials, time.Hour*8760); err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateRequest")
	}
	tests.Passed("Should have generated new CertificateRequest")

	raw, err := reqCA.Raw()
	if err != nil {
		tests.FailedWithError(err, "Should have generated raw version of CertificateRequest")
	}
	tests.Passed("Should have generated raw version of CertificateRequest")

	var rca certificates.CertificateRequest
	if err := rca.FromRaw(raw); err != nil {
		tests.FailedWithError(err, "Should have read raw  of CertificateRequest")
	}
	tests.Passed("Should have read raw  of CertificateRequest")
}

func TestCertificateService(t *testing.T) {
	storeMap := make(map[string][]byte)
	var store mocks.PersistenceStoreMock
	store.GetFunc, store.PersistFunc = mocks.MapStore(storeMap)

	serials := certificates.SerialService{Length: 128}
	profile := certificates.CertificateProfile{
		Local:        "Lagos",
		Organization: "DreamBench",
		CommonName:   "DreamBench Inc",
		Country:      "Nigeria",
		Province:     "South-West",
	}

	var service certificates.CertificateAuthorityService
	service.KeyStrength = 4096
	service.LifeTime = (time.Hour * 8760)
	service.Profile = profile
	service.Serials = serials
	service.Emails = append([]string{}, "alex.ewetumo@dreambench.io")

	ca, err := service.New()
	if err != nil {
		tests.FailedWithError(err, "Should have generated new CertificateAuthority")
	}
	tests.Passed("Should have generated new CertificateAuthority")

	if err := ca.Persist(store); err != nil {
		tests.FailedWithError(err, "Should have successfully store certificate into persistence store")
	}
	tests.Passed("Should have successfully store certificate into persistence store")

	var restoredCA certificates.CertificateAuthority
	if err := restoredCA.Load(store); err != nil {
		tests.FailedWithError(err, "Should have successfully retrieved certificate from store")
	}
	tests.Passed("Should have successfully retrieved certificate from store")

	caRaw, err := ca.CertificateRaw()
	if err != nil {
		tests.FailedWithError(err, "Should have been able to retrieve raw form of certificate")
	}
	tests.Passed("Should have been able to retrieve raw form of certificate")

	rcaRaw, err := restoredCA.CertificateRaw()
	if err != nil {
		tests.FailedWithError(err, "Should have been able to retrieve raw form of certificate")
	}
	tests.Passed("Should have been able to retrieve raw form of certificate")

	if !bytes.Equal(caRaw, rcaRaw) {
		tests.Failed("Should have matching certificate raw data between real and restored versions")
	}
	tests.Passed("Should have matching certificate raw data between real and restored versions")

	caKeyRaw, err := ca.PrivateKeyRaw()
	if err != nil {
		tests.FailedWithError(err, "Should have been able to retrieve raw form of certificate")
	}
	tests.Passed("Should have been able to retrieve raw form of certificate")

	rcaKeyRaw, err := restoredCA.PrivateKeyRaw()
	if err != nil {
		tests.FailedWithError(err, "Should have been able to retrieve raw form of certificate")
	}
	tests.Passed("Should have been able to retrieve raw form of certificate")

	if !bytes.Equal(caKeyRaw, rcaKeyRaw) {
		tests.Failed("Should have matching certificate private key raw data between real and restored versions")
	}
	tests.Passed("Should have matching certificate private key raw data between real and restored versions")

	caConRaw, err := ca.Raw()
	if err != nil {
		tests.FailedWithError(err, "Should have being able to generate raw bytes of CertificateAuthority")
	}
	tests.Passed("Should have being able to generate raw bytes of CertificateAuthority")

	var newCA certificates.CertificateAuthority
	if err := newCA.FromRaw(caConRaw); err != nil {
		tests.FailedWithError(err, "Should have being able to load raw version of CertificateAuthority")
	}
	tests.Passed("Should have being able to load raw version of CertificateAuthority")

	rcaRaw, err = newCA.CertificateRaw()
	if err != nil {
		tests.FailedWithError(err, "Should have been able to retrieve raw form of certificate")
	}
	tests.Passed("Should have been able to retrieve raw form of certificate")

	if !bytes.Equal(caRaw, rcaRaw) {
		tests.Failed("Should have matching certificate raw data between real and restored versions")
	}
	tests.Passed("Should have matching certificate raw data between real and restored versions")

	rcaKeyRaw, err = newCA.PrivateKeyRaw()
	if err != nil {
		tests.FailedWithError(err, "Should have been able to retrieve raw form of certificate")
	}
	tests.Passed("Should have been able to retrieve raw form of certificate")

	if !bytes.Equal(caKeyRaw, rcaKeyRaw) {
		tests.Failed("Should have matching certificate private key raw data between real and restored versions")
	}
	tests.Passed("Should have matching certificate private key raw data between real and restored versions")
}
