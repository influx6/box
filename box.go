package box

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/BurntSushi/toml"
)

// variables.
var (
	functions = map[string]Generator{}
)

// errors.
var (
	ErrFunctionNotFound = errors.New("Function with given id not found")
)

// CancelContext defines a type which provides Done signal for cancelling operations.
type CancelContext interface {
	Done() <-chan struct{}
}

// Spell defines an interface which expose an exec method.
type Spell interface {
	Exec(CancelContext) error
}

// Function defines a interface for a function which returns a giving Spell.
type Function func() Spell

// Generator defines a interface for a function which returns a giving Spell.
type Generator func([]byte) (Spell, error)

// Register adds the giving function into the list of available functions for instanction.
func Register(id string, fun Generator) bool {
	functions[id] = fun
	return true
}

// RegisterTOML adds the giving function into the list of available functions for instanction,
// where it's config would be loaded using toml as the config unmarshaller.
func RegisterTOML(id string, fun Function) bool {
	functions[id] = func(config []byte) (Spell, error) {
		instance := fun()
		if _, err := toml.DecodeReader(bytes.NewBuffer(config), instance); err != nil {
			return nil, err
		}

		return instance, nil
	}

	return true
}

// RegisterJSON adds the giving function into the list of available functions for instanction,
// where it's config would be loaded using json as the config unmarshaller.
func RegisterJSON(id string, fun Function) bool {
	functions[id] = func(config []byte) (Spell, error) {
		instance := fun()
		if err := json.Unmarshal(config, instance); err != nil {
			return nil, err
		}

		return instance, nil
	}

	return true
}

// MustCreateFromBytes panics if giving function for id is not found.
func MustCreateFromBytes(id string, config []byte) Spell {
	spell, err := CreateFromBytes(id, config)
	if err != nil {
		panic(err)
	}

	return spell
}

// CreateFromBytes returns a new spell from the provided configuration
func CreateFromBytes(id string, config []byte) (Spell, error) {
	fun, ok := functions[id]
	if !ok {
		return nil, ErrFunctionNotFound
	}

	return fun(config)
}

// CreateWithTOML returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func CreateWithTOML(id string, config map[string]interface{}) (Spell, error) {
	var bu bytes.Buffer

	if err := toml.NewEncoder(&bu).Encode(config); err != nil {
		return nil, err
	}

	return CreateFromBytes(id, bu.Bytes())
}

// CreateWithJSON returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func CreateWithJSON(id string, config map[string]interface{}) (Spell, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return CreateFromBytes(id, data)
}
