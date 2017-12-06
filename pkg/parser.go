package pkg

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

const (
	maxOPNameLength = 20
)

// errs ...
var (
	ErrNoIdentity     = errors.New("message must have identity info")
	ErrNoOpName       = errors.New("message must have op name after secret identity hash")
	ErrNoSecretHeader = errors.New("message must have secret with hash suffix")
	ErrInvalidSecret  = errors.New("message secret does not match expect secret")
	ErrNoHash         = errors.New("message must have initial hash with secret prefix")
	ErrInvalidOpName  = errors.New("message must have op name with max length of 50 in characters")
	ErrInvalidMessage = errors.New("messages must be in format: secret#identity OP [body...]")
)

var (
	spaceRune = rune(' ')
	hashRune  = rune('#')
)

// Op defines a structure which is used to contain data
// of message.
type Op struct {
	Op       string `json:"op"`
	Secret   string `json:"secret"`
	Identity string `json:"identity"`
	Body     []byte `json:"body"`
	Raw      []byte `json:"raw"`
}

// OpParser implements a message parser which request messages to occur
// in the giving below format. Which it breaks down into appropriate
// structure which will contain necessary information from origin.
// Format: secret#identity OP [....]
type OpParser struct {
	Secret string
}

// Parse implements the necessary logic needed for creating a appropriate Op
// structure which would contain the necessary bits else returning an error
// if the data is invalid.
func (op OpParser) Parse(incoming []byte) (Op, error) {
	var msg Op
	msg.Raw = incoming

	reader := bufio.NewReader(bytes.NewReader(incoming))
	prosecret, err := op.pullSecret(reader)
	if err != nil && err != io.EOF {
		return msg, err
	}

	if err != nil && err == io.EOF {
		return msg, ErrInvalidMessage
	}

	if string(prosecret) != op.Secret {
		return msg, ErrInvalidSecret
	}

	identity, err := op.pullIdentity(reader)
	if err != nil && err != io.EOF {
		return msg, err
	}

	if err != nil && err == io.EOF {
		return msg, ErrInvalidMessage
	}

	msg.Identity = string(identity)

	opHeader, err := op.pullOPName(reader)
	if err != nil && err != io.EOF {
		return msg, err
	}

	msg.Op = string(opHeader)

	body, err := op.pullBody(reader)
	if err != nil {
		return msg, err
	}

	msg.Body = body

	return msg, nil
}

func (op OpParser) pullBody(reader *bufio.Reader) ([]byte, error) {
	var started bool

	// create space for data with enough space.
	body := make([]byte, 0, 1024)
	for {
		bit, err := reader.ReadByte()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err != nil && err == io.EOF {
			return body, nil
		}

		switch rune(bit) {
		case spaceRune:
			if !started {
				continue
			}
			fallthrough
		default:
			started = true
			body = append(body, bit)
		}
	}
}

func (op OpParser) pullOPName(reader *bufio.Reader) ([]byte, error) {
	var started bool
	// create space for data with enough space.
	opName := make([]byte, 0, maxOPNameLength)

pullLoop:
	for {
		bit, err := reader.ReadByte()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err != nil && err == io.EOF {
			if len(opName) == 0 {
				return nil, ErrNoOpName
			}

			if len(opName) >= maxOPNameLength {
				return nil, ErrInvalidOpName
			}

			return opName, err
		}

		if started && len(opName) >= maxOPNameLength {
			return nil, ErrInvalidOpName
		}

		switch rune(bit) {
		case spaceRune:
			if started {
				break pullLoop
			}
		default:
			started = true
			opName = append(opName, bit)
		}
	}

	if len(opName) >= maxOPNameLength {
		return nil, ErrInvalidOpName
	}

	return opName, nil
}

func (op OpParser) pullIdentity(reader *bufio.Reader) ([]byte, error) {
	// create space for data with enough space.
	identity := make([]byte, 0, 128)

pullLoop:
	for {
		bit, err := reader.ReadByte()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err != nil && err == io.EOF {
			if len(identity) == 0 {
				return nil, ErrNoIdentity
			}
			return identity, err
		}

		switch rune(bit) {
		case spaceRune:
			break pullLoop
		default:
			identity = append(identity, bit)
		}
	}
	return identity, nil
}

func (op OpParser) pullSecret(reader *bufio.Reader) ([]byte, error) {
	// create space with enough space.
	secret := make([]byte, 0, 128)
pullLoop:
	for {
		bit, err := reader.ReadByte()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err != nil && err == io.EOF {
			if len(secret) == 0 {
				return nil, ErrSecretRequired
			}
			return secret, err
		}

		switch rune(bit) {
		case hashRune:
			break pullLoop
		case spaceRune:
			if len(secret) != 0 {
				return nil, ErrNoSecretHeader
			}
		default:
			secret = append(secret, bit)
		}
	}
	return secret, nil
}
