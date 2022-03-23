package httputils

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	InnerMock mock.Mock
}

func (mock *MockHTTPClient) Do(request *http.Request) (*http.Response, error) {
	results := mock.InnerMock.Called(request)
	return results.Get(0).(*http.Response), results.Error(1)
}
