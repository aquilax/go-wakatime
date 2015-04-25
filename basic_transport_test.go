package wakatime

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestBasicTransport(t *testing.T) {
	Convey("Given basic transport", t, func() {
		bt := NewBasicTransport("key")
		Convey("The transport must not be nil", func() {
			So(bt, ShouldNotBeNil)
		})
		Convey("Request must have the correct headers", func() {
			r, err := http.NewRequest("GET", "http://example.com/", nil)
			So(err, ShouldBeNil)
			req := bt.addHeaders(r)
			Convey("Authorization must be correct", func() {
				So(req.Header.Get("Authorization"), ShouldEqual, "Basic a2V5")
			})
			Convey("Content type must be correct", func() {
				So(req.Header.Get("Content-Type"), ShouldEqual, "application/json")
			})
			Convey("User agent  must be correct", func() {
				So(req.Header.Get("User-Agent"), ShouldEqual, getUserAgent())
			})
		})
	})
}
