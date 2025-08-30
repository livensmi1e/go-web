package mocks

import "github.com/stretchr/testify/mock"

type MockToken struct {
	mock.Mock
}

func (m *MockToken) Generate(claims map[string]interface{}) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}

func (m *MockToken) Validate(token string) (map[string]interface{}, error) {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
