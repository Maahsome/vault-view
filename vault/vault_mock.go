package vault

import (
	"github.com/stretchr/testify/mock"
)

type vaultMock struct {
	mock.Mock
}

// NewVaultMock - Mocking the vault interactions
func NewVaultMock() Client {

	return &vaultMock{}

}

func (vm *vaultMock) GetVersion() (string, error) {

	return "1.2.3", nil
}

func (vm *vaultMock) GetData(keypath string) (DataRecord, error) {

	return DataRecord{}, nil

}
func (vm *vaultMock) GetPaths(keypath string) (map[string]Paths, error) {

	return map[string]Paths{}, nil
}
