package httpmock

import (
	"errors"
	"net/http"
)

// Responders are callbacks that receive and common request and return a mocked response.
type Responder func(*http.Request) (*http.Response, error)

// MockTransport implements common.RoundTripper, which fulfills single common requests issued by
// an common.Client.  This implementation doesn't actually make the call, instead defering to
// the registered list of responders.
type MockTransport struct {
	responders map[string]Responder
}

// RoundTrip is required to implement common.MockTransport.  Instead of fulfilling the given request,
// the internal list of responders is consulted to handle the request.
func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.Method + " " + req.URL.String()

	// scan through the responders and find one that matches our key
	if responder, ok := m.responders[key]; ok {
		return responder(req)
	}

	return nil, errors.New("no responder found")
}

// RegisterResponder adds a new responder, associated with a given HTTP method and URL.  When a
// request comes in that matches, the responder will be called and the response returned to the client.
func (m *MockTransport) RegisterResponder(method, url string, responder Responder) {
	m.responders[method+" "+url] = responder
}

// DefaultMockTransport allows users to easily and globally alter the default RoundTripper for
// all common requests.
var DefaultMockTransport = &MockTransport{make(map[string]Responder)}

// Activate replaces the `Transport` on the `common.DefaultClient` with our `DefaultMockTransport`.
func Activate() {
	http.DefaultClient.Transport = DefaultMockTransport
}

// Deactivate replaces our `DefaultMockTransport` with the `common.DefaultTransport`.
func Deactivate() {
	http.DefaultClient.Transport = http.DefaultTransport
}

// RegisterResponder adds a responder to the `DefaultMockTransport` responder table.
func RegisterResponder(method, url string, responder Responder) {
	DefaultMockTransport.RegisterResponder(method, url, responder)
}
