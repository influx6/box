package box

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/influx6/faux/context"
)

// variables.
var (
	functions = NewGeneratorRegistry()
)

// errors.
var (
	ErrFunctionNotFound = errors.New("Function with given id not found")
)

// Spell defines an interface which expose an exec method.
type Spell interface {
	Exec(context.CancelContext) error
}

// Function defines a interface for a function which returns a giving Spell.
type Function func() Spell

// Generator defines a interface for a function which returns a giving Spell.
type Generator func([]byte) (Spell, error)

// Register adds the giving function into the list of available functions for instanction.
func Register(id string, fun Generator) bool {
	return functions.Register(id, fun)
}

// RegisterTOML adds the giving function into the list of available functions for instanction,
// where it's config would be loaded using toml as the config unmarshaller.
func RegisterTOML(id string, fun Function) bool {
	return functions.RegisterTOML(id, fun)
}

// RegisterJSON adds the giving function into the list of available functions for instanction,
// where it's config would be loaded using json as the config unmarshaller.
func RegisterJSON(id string, fun Function) bool {
	return functions.RegisterJSON(id, fun)
}

// MustCreateFromBytes panics if giving function for id is not found.
func MustCreateFromBytes(id string, config []byte) Spell {
	return functions.MustCreateFromBytes(id, config)
}

// CreateFromBytes returns a new spell from the provided configuration
func CreateFromBytes(id string, config []byte) (Spell, error) {
	return functions.CreateFromBytes(id, config)
}

// CreateWithTOML returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func CreateWithTOML(id string, config map[string]interface{}) (Spell, error) {
	return functions.CreateWithTOML(id, config)
}

// CreateWithJSON returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func CreateWithJSON(id string, config map[string]interface{}) (Spell, error) {
	return functions.CreateWithJSON(id, config)
}

// GeneratorRegistry implements a structure that handles registration and ochestration of spells.
type GeneratorRegistry struct {
	functions map[string]Generator
}

// NewGeneratorRegistry returns a new instance of the GeneratorsRegistry.
func NewGeneratorRegistry() *GeneratorRegistry {
	return &GeneratorRegistry{
		functions: make(map[string]Generator),
	}
}

// Register adds the giving function into the list of available functions for instanction.
func (ng *GeneratorRegistry) Register(id string, fun Generator) bool {
	ng.functions[id] = fun
	return true
}

// RegisterTOML adds the giving function into the list of available functions for instanction,
// where it's config would be loaded using toml as the config unmarshaller.
func (ng *GeneratorRegistry) RegisterTOML(id string, fun Function) bool {
	ng.functions[id] = func(config []byte) (Spell, error) {
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
func (ng *GeneratorRegistry) RegisterJSON(id string, fun Function) bool {
	ng.functions[id] = func(config []byte) (Spell, error) {
		instance := fun()
		if err := json.Unmarshal(config, instance); err != nil {
			return nil, err
		}

		return instance, nil
	}

	return true
}

// MustCreateFromBytes panics if giving function for id is not found.
func (ng *GeneratorRegistry) MustCreateFromBytes(id string, config []byte) Spell {
	spell, err := ng.CreateFromBytes(id, config)
	if err != nil {
		panic(err)
	}

	return spell
}

// CreateFromBytes returns a new spell from the provided configuration
func (ng *GeneratorRegistry) CreateFromBytes(id string, config []byte) (Spell, error) {
	fun, ok := ng.functions[id]
	if !ok {
		return nil, ErrFunctionNotFound
	}

	return fun(config)
}

// CreateWithTOML returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func (ng *GeneratorRegistry) CreateWithTOML(id string, config map[string]interface{}) (Spell, error) {
	var bu bytes.Buffer

	if err := toml.NewEncoder(&bu).Encode(config); err != nil {
		return nil, err
	}

	return ng.CreateFromBytes(id, bu.Bytes())
}

// CreateWithJSON returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func (ng *GeneratorRegistry) CreateWithJSON(id string, config map[string]interface{}) (Spell, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return ng.CreateFromBytes(id, data)
}
