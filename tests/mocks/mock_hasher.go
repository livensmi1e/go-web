package mocks

import "github.com/stretchr/testify/mock"

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) Compare(hash string, plain string) error {
	args := m.Called(hash, plain)
	return args.Error(0)
}
