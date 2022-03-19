package httputils

import "github.com/stretchr/testify/mock"

type MockLogger struct {
	InnerMock mock.Mock
}

func (mock *MockLogger) Print(arguments ...interface{}) {
	mock.InnerMock.Called(arguments)
}
