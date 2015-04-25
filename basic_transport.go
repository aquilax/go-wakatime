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
		encodedAPIKey: base64.StdEncoding.EncodeToString([]byte(apiKey)),
		Transport:     http.DefaultTransport,
	}
}

func (bt *BasicTransport) addHeaders(req *http.Request) *http.Request {
	req = cloneRequest(req)
	req.Header.Set("Authorization", "Basic "+bt.encodedAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", getUserAgent())
	return req
}

// RoundTrip implements the http.RoundTripper method
func (bt *BasicTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return bt.Transport.RoundTrip(bt.addHeaders(req))
}

func getUserAgent() string {
	return "go-wakatime/" + Version
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
