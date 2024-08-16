package storage

import "github.com/stretchr/testify/mock"

// MockStorage is a mock implementation of Storage for testing.
type MockStorage struct {
	mock.Mock
}

// GetWallet mocks the GetWallet method of Storage.
func (m *MockStorage) GetWallet() (*Wallet, error) {
	args := m.Called()
	return args.Get(0).(*Wallet), args.Error(1)
}
