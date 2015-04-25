package wakatime

import (
	"encoding/base64"
	"net/http"
)

// BasicTransport implements http.RoundTripper and provides authentication using
// the Basic Auth mechanism
type BasicTransport struct {
	encodedAPIKey string
	Transport     http.RoundTripper
}

// NewBasicTransport creates new BasicTransport given the API key
func NewBasicTransport(apiKey string) *BasicTransport {
	return &BasicTransport{
		encodedApiKey: base64.StdEncoding.EncodeToString([]byte(apiKey)),
	}
}

// RoundTrip implements the http.RoundTripper method
func (bt *BasicTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req)
	req.Header.Set("Authorization", "Basic "+bt.encodedAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go-wakatime/"+Version)

	// Make the HTTP request.
	return bt.transport().RoundTrip(req)
}

func (bt *BasicTransport) transport() http.RoundTripper {
	if bt.Transport != nil {
		return bt.Transport
	}
	return http.DefaultTransport
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}
