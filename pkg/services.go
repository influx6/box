package pkg

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/influx6/faux/metrics"
	"github.com/influx6/mnet"
	"github.com/influx6/mnet/certificates"
)

// SecurityCtrl implements the certificate authorization and request system which
// is used to initiate a secure communication with the ContainerCtrl through the network.
type SecurityCtrl struct {
	Box     *BoxCtrl
	Metrics metrics.Metrics
}

// Serve implements the necessary method which receives messages over the wire
// through provided mtcp.Client with appropriate response.
//
// SecurityCtrl supports the following operations:
//
// 	Please remember every message must have 'secret#identity' prefix, we use [IDENTITY] as placeholder.
//	If associated secret and identity is banned or invalid, any request will be rejected immediately.
//
//	[IDENTITIY] REGCL [BASE64 ENCODED CERTIFICATE REQUEST BINARY]
//		REGCL protocol seeks registeration of client as a secured client, by issueing it's information
//		and attaching a certificate request seeking the server to return a signed
//		certificate which it can use to securly connect to issue orders.
//
//		REGCL responds in two ways:
//		1. Success		[SERVER_IDENTITY] REGCLRES +OK [CERITIFCATE RESPONSE]
//		2. Failure		[SERVER_IDENTITY] REGCLRES +ERR [CERITIFCATE ERROR]
//
//		Where [SERVER_IDENTITY] is secret#server_identity.
//
//
func (sec SecurityCtrl) Serve(client *mnet.Client) error {
	// create parser for message processing.
	parser := OpParser{Secret: sec.Box.Secret}

	for {
		message, err := client.Read()
		if err != nil {
			if err == mnet.ErrNoDataYet {
				continue
			}
			sec.Metrics.Emit(
				metrics.Error(err),
				metrics.With("box", sec.Box.ID),
				metrics.With("controller", sec.Box.ID),
				metrics.Message("Failed reading mnet.Client"),
			)
			return err
		}

		msg, err := parser.Parse(message)
		if err != nil {
			sec.Metrics.Emit(
				metrics.Error(err),
				metrics.With("box", sec.Box.ID),
				metrics.With("controller", sec.Box.ID),
				metrics.Message("Failed reading mnet.Client"),
			)
			return err
		}

		switch msg.Op {
		case "REGCL":
			return sec.registerClient(msg, client)
		default:
			// Nothing to do here.
		}
	}
}

// registerClient implements the logic to process incoming message and
// generate appropriate files for new client. Any previous data of client
// is wiped out, recreating once again.
func (sec SecurityCtrl) registerClient(op Op, client *mnet.Client) error {

	// Convert client data from base64 to json format.
	clientCertificateRequest, err := base64.StdEncoding.Decode(op.Body)
	if err != nil {
		return err
	}

	req, err := x509.ParseCertificateRequest(clientCertificateRequest)
	if err != nil {
		return err
	}

	var certReq certificates.CertificateRequest
	certReq.Request = req

	if err := sec.Box.signCertificateRequest(&certReq); err != nil {
		fmt.Fprintf(client, "%s#%s REGCLRES +ERR %+s", sec.Box.Secret, sec.Box.ID, err.Error())
		return client.Flush()
	}

	rawResponse, err := certReq.Raw()
	if err != nil {
		fmt.Fprintf(client, "%s#%s REGCLRES +ERR %+s", sec.Box.Secret, sec.Box.ID, err.Error())
		return client.Flush()
	}

	if _, err := fmt.Fprintf(client, "%s#%s REGCLRES +OK %+s", sec.Box.Secret, sec.Box.ID, base64.StdEncoding.EncodeToString(rawResponse)); err != nil {
		return err
	}

	return client.Flush()
}
