package httpmock

import (
	"errors"
	"net/http"
)

type Responder func(*http.Request) (*http.Response, error)

type MockTransport struct {
	responders map[string]Responder
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.Method + " " + req.URL.String()

	if responder, ok := m.responders[key]; ok {
		return responder(req)
	}

	return nil, errors.New("no responder found")
}

func (m *MockTransport) RegisterResponder(method, url string, responder Responder) {
	m.responders[method+" "+url] = responder
}

var DefaultMockTransport = &MockTransport{make(map[string]Responder)}

func Activate() {
	http.DefaultClient.Transport = DefaultMockTransport
}

func Deactivate() {
	http.DefaultClient.Transport = http.DefaultTransport
}

func RegisterResponder(method, url string, responder Responder) {
	DefaultMockTransport.RegisterResponder(method, url, responder)
}
