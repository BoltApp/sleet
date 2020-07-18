package testing

import (
	"net/http"
	"net/url"
	"testing"
)

// MockHttpClient provides the ability to add testing functionality ontop of a standard http client
type MockHttpClient struct {
	Client        *http.Client
	TestServerURL string
	RealURL       string
	T             *testing.T
}

// Do checks if the URL set on the http Request matches an expected value
func (m MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	m.T.Helper()

	if req.URL.String() != m.RealURL {
		m.T.Errorf("Wrong URL called, Got %q, want %q", req.URL.String(), m.RealURL)
	}

	req.URL, _ = url.Parse(m.TestServerURL)

	return m.Client.Do(req)
}
