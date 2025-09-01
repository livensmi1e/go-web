package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Set(key string, value interface{}) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockCache) Get(key string, value interface{}) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockCache) SetWithTTL(key string, value interface{}, ttl int) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}
