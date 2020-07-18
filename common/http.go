package common

import (
	"fmt"
	"net/http"
	"time"
)

// DefaultTimeout for HTTP Requests
const DefaultTimeout = 60 * time.Second

// DefaultHttpClient returns an http client with defaultTimeout specified above
func DefaultHttpClient() *http.Client {
	return &http.Client{
		Timeout: DefaultTimeout,
	}
}

type HttpSender interface {
	Do(*http.Request) (*http.Response, error)
}

// UserAgent specifies the Sleet library and version for PsPs that require this header
func UserAgent() string {
	return fmt.Sprintf("Sleet/%s", LibraryVersion)
}
