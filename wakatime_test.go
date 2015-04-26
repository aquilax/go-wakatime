package wakatime

import (
	"bufio"
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	users = `{
  "data": {
    "created": "2015-04-23T04:32:26Z",
    "email": "aquilax@example.com",
    "email_public": true,
    "full_name": "Full Name",
    "human_readable_website": "www.avtobiografia.com",
    "id": "e9b45851-991b-4755-ffff-6355d927f472",
    "last_heartbeat": "2015-04-26T04:16:23Z",
    "last_plugin": "wakatime/4.0.8 (Linux-3.16.0-34-generic-x86_64-with-Ubuntu-14.10-utopic) Python2.7.8.final.0 sublime/3083 sublime-wakatime/4.0.0",
    "last_plugin_name": "Sublime",
    "last_project": "go-wakatime",
    "location": "Stockholm, Sweden",
    "logged_time_public": false,
    "modified": "2015-04-23T05:48:57Z",
    "photo": "https://secure.gravatar.com/avatar/8cf592a20de754300721bf954aa40507?s=150&d=identicon",
    "photo_public": true,
    "plan": "basic",
    "timezone": "Europe/Stockholm",
    "username": "aquilax",
    "website": "http://www.avtobiografia.com"
  }
}`
)

type DummyTransport struct {
	content string
}

func NewDummyTransport(content string) *DummyTransport {
	return &DummyTransport{content}
}

func (dt *DummyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	b := bytes.NewBufferString(dt.content)
	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bufio.NewReader(b)),
	}
	return resp, nil //http.ReadResponse(bufio.NewReader(b), nil)
}

func TestWakatime(t *testing.T) {
	Convey("Given wakatime", t, func() {
		wt := New(NewDummyTransport(users))
		Convey("Wakatime must not be nil", func() {
			So(wt, ShouldNotBeNil)
			Convey("JSON must be correctly parsed", func() {
				u, err := wt.Users(CurrentUser)
				So(err, ShouldBeNil)
				So(u, ShouldNotBeNil)
				So(u.Data, ShouldNotBeNil)
				So(u.Data.Created.Format(time.RFC3339), ShouldEqual, "2015-04-23T04:32:26Z")
				So(u.Data.Email, ShouldEqual, "aquilax@example.com")
				So(u.Data.EmailPublic, ShouldBeTrue)
				So(u.Data.FullName, ShouldEqual, "Full Name")
				So(u.Data.HumanReadableWebsite, ShouldEqual, "www.avtobiografia.com")
				So(u.Data.ID, ShouldEqual, "e9b45851-991b-4755-ffff-6355d927f472")
			})
		})
	})
}
