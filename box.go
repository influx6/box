package box

import (
	"github.com/influx6/faux/ops"
)

// variables.
var (
	functions = ops.NewGeneratorRegistry()
)

// Register adds the giving function into the list of available functions for instanction.
func Register(id string, fun ops.Generator) bool {
	return functions.Register(id, fun)
}

// RegisterTOML adds the giving function into the list of available functions for instanction,
// where it's config would be loaded using toml as the config unmarshaller.
func RegisterTOML(id string, fun ops.Function) bool {
	return functions.RegisterTOML(id, fun)
}

// RegisterJSON adds the giving function into the list of available functions for instanction,
// where it's config would be loaded using json as the config unmarshaller.
func RegisterJSON(id string, fun ops.Function) bool {
	return functions.RegisterJSON(id, fun)
}

// MustCreateFromBytes panics if giving function for id is not found.
func MustCreateFromBytes(id string, config []byte) ops.Op {
	return functions.MustCreateFromBytes(id, config)
}

// CreateFromBytes returns a new spell from the provided configuration
func CreateFromBytes(id string, config []byte) (ops.Op, error) {
	return functions.CreateFromBytes(id, config)
}

// CreateWithTOML returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func CreateWithTOML(id string, config map[string]interface{}) (ops.Op, error) {
	return functions.CreateWithTOML(id, config)
}

// CreateWithJSON returns a new spell from the provided configuration map which is first
// converted into JSON then loaded using the CreateFromBytes function.
func CreateWithJSON(id string, config map[string]interface{}) (ops.Op, error) {
	return functions.CreateWithJSON(id, config)
}
