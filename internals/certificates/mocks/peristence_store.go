package mocks

import "errors"

// PersistenceStoreMock implements a mock structure for the
// certificates.PersistenceStore interface.
type PersistenceStoreMock struct {
	GetFunc     func(string) ([]byte, error)
	PersistFunc func(string, []byte) error
}

// Retrieve attempts to retrieve giving item from the store.
func (p PersistenceStoreMock) Retrieve(name string) ([]byte, error) {
	if p.GetFunc == nil {
		return nil, errors.New("Unable to get name")
	}

	return p.GetFunc(name)
}

// Persist attempts adding giving data into store keyed by provided name.
func (p PersistenceStoreMock) Persist(name string, data []byte) error {
	if p.PersistFunc == nil {
		return errors.New("Unable to get name")
	}

	return p.PersistFunc(name, data)
}

// MapStore returns two functions which store and retrieve items from provided map.
func MapStore(store map[string][]byte) (getFunc func(string) ([]byte, error), persistFunc func(string, []byte) error) {
	getFunc = func(name string) ([]byte, error) {
		if val, ok := store[name]; ok {
			return val, nil
		}
		return nil, errors.New("not found")
	}

	persistFunc = func(name string, data []byte) error {
		store[name] = data
		return nil
	}
	return
}
