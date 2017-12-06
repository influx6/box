package pkg

import (
	"github.com/influx6/mnet"
	"github.com/influx6/mnet"
)

// SecurityCtrl implements the certificate authorization and request system which
// is used to initiate a secure communication with the ContainerCtrl through the network.
type SecurityCtrl struct {
	box *BoxCtrl
}

// Serve implements the necessary method which receives messages over the wire
// through provided mtcp.Client with appropriate response.
func (sec SecurityCtrl) Serve(client *mnet.Client) error {
	var err error
	var message []byte

	for {
		message, err = client.Read()
		if err != nil {
			if err == mnet.ErrNoDataYet {
				err = nil
				continue
			}
			break
		}

		_ = message
	}

	return err
}
