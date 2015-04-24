package wakatime

import (
	"net/http"
)

type BasicTransport struct {
	apiKey    string
	Transport http.RoundTripper
}

func NewBasicTransport(apiKey string) *BasicTransport {
	return &BasicTransport{
		apiKey: apiKey,
	}
}

func (bt *BasicTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req)
	req.Header.Set("Authorization", "Basic "+bt.apiKey)

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
