package wakatime

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

type DummyTransport struct{}

func NewDummyTransport() *DummyTransport {
	return &DummyTransport{}
}

func (dt *DummyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, nil
}

func TestWakatime(t *testing.T) {
	Convey("Given wakatime", t, func() {
		wt := New(NewDummyTransport())
		Convey("Wakatime must not be nil", func() {
			So(wt, ShouldNotBeNil)
		})
	})
}
