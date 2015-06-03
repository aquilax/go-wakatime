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
    "last_plugin": "wakatime/4.0.8",
    "last_plugin_name": "Sublime",
    "last_project": "go-wakatime",
    "location": "Stockholm, Sweden",
    "logged_time_public": true,
    "modified": "2015-04-23T05:48:57Z",
    "photo": "https://secure.gravatar.com/avatar/8cf592a20de754300721bf954aa40507?s=150&d=identicon",
    "photo_public": true,
    "plan": "basic",
    "timezone": "Europe/Stockholm",
    "username": "aquilax",
    "website": "http://www.avtobiografia.com"
  }
}`
	summaries = `{
  "data": [
    {
      "editors": [
        {
          "digital": "2:1005",
          "hours": 2,
          "minutes": 10,
          "name": "PhpStorm",
          "percent": 69.91,
          "seconds": 5,
          "text": "2 hours 10 minutes 5 seconds",
          "total_seconds": 7805
        }
      ],
      "grand_total": {
        "digital": "3:03",
        "hours": 3,
        "minutes": 3,
        "text": "3 hours 3 minutes",
        "total_seconds": 11165
      },
      "languages": [
        {
          "digital": "0:22:58",
          "hours": 0,
          "minutes": 22,
          "name": "Go",
          "percent": 12.34,
          "seconds": 58,
          "text": "22 minutes 58 seconds",
          "total_seconds": 1378
        }
      ],
      "operating_systems": [
        {
          "digital": "3:0604",
          "hours": 3,
          "minutes": 6,
          "name": "Linux",
          "percent": 100,
          "seconds": 4,
          "text": "3 hours 6 minutes 4 seconds",
          "total_seconds": 11164
        }
      ],
      "projects": [
        {
          "digital": "0:22",
          "hours": 0,
          "minutes": 22,
          "name": "go-wakatime",
          "percent": 12.23,
          "text": "22 minutes",
          "total_seconds": 1365
        }
      ],
      "range": {
        "date": "04\/23\/2015",
        "date_human": "04\/23\/2015",
        "end": 1429826399,
        "start": 1429740000,
        "text": "04\/23\/2015",
        "timezone": "Europe\/Stockholm"
      }
    }
  ],
  "end": 1429912799,
  "start": 1429740000
}`
	durations = `{
  "branches": [
    "master"
  ],
  "data": [
    {
      "duration": 2240.0,
      "project": "go-wakatime",
      "time": 1430021746.422815
    }
  ],
  "end": 1430085599,
  "start": 1429999200,
  "timezone": "Europe/Stockholm"
}`
	stats = `{
  "data": {
    "created_at": "2015-04-23T04:38:05Z",
    "editors": [
      {
        "created_at": "2015-04-24T07:12:29Z",
        "id": "a5979e20-aa2a-403d-9214-b49e85f00fbc",
        "modified_at": "2015-04-28T07:08:06Z",
        "name": "Vim",
        "percent": 22.64,
        "total_seconds": 11803
      }
    ],
    "end": 1430171999.000000,
    "human_readable_daily_average": "2 hours 3 minutes",
    "human_readable_total": "14 hours 24 minutes",
    "id": "3e570b91-2540-4c9e-a71a-75b1909188ea",
    "is_up_to_date": true,
    "languages": [
      {
        "created_at": "2015-04-24T07:12:29Z",
        "id": "23c0dd3f-d09a-4c01-a813-37aa257a114c",
        "modified_at": "2015-04-28T07:08:06Z",
        "name": "Go",
        "percent": 41.37,
        "total_seconds": 21569
      }
    ],
    "modified_at": "2015-04-28T07:08:06Z",
    "operating_systems": [
      {
        "created_at": "2015-04-24T07:12:29Z",
        "id": "36685edf-5a5e-4d53-8d77-61c1dd6e9c46",
        "modified_at": "2015-04-28T07:08:06Z",
        "name": "Linux",
        "percent": 100.00,
        "total_seconds": 52137
      }
    ],
    "project": null,
    "projects": [
      {
        "created_at": "2015-04-24T07:12:29Z",
        "id": "198e4ed6-a208-41b7-b698-1826a003411a",
        "modified_at": "2015-04-28T07:08:06Z",
        "name": "go-wakatime",
        "percent": 45.77,
        "total_seconds": 23865
      }
    ],
    "range": "last_7_days",
    "start": 1429567200.000000,
    "status": "ok",
    "timeout": 15,
    "timezone": "Europe/Stockholm",
    "total_seconds": 51840,
    "user_id": "e9b45851-991b-4755-9ccd-6355d927f472",
    "username": "aquilax",
    "writes_only": true
  }
}`
	hartbeat = `{
  "data": [
    {
      "branch": "master",
      "entity": "/home/aquilax/projects/frondate_project/frondate/settings/production.py",
      "id": "2727163e-e614-47cd-8ffc-f3dc3f18bcdc",
      "is_debugging": null,
      "is_write": false,
      "language": "Python",
      "project": "frondate_project",
      "time": 1433217822.482732,
      "type": "file"
    },
    {
      "branch": "master",
      "entity": "/home/aquilax/projects/frondate_project/vagrant/debian-jessie/provision.sh",
      "id": "690c5596-33f1-458a-9a8f-8d4062daa5bb",
      "is_debugging": null,
      "is_write": false,
      "language": "Bash",
      "project": "frondate_project",
      "time": 1433217827.417102,
      "type": "file"
    }
  ],
  "end": 1433282399,
  "start": 1433196000,
  "timezone": "Europe/Stockholm"
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
	return resp, nil
}

func TestWakatime(t *testing.T) {
	Convey("Given wakatime", t, func() {
		wt := New(NewDummyTransport(users))
		Convey("Wakatime must not be nil", func() {
			So(wt, ShouldNotBeNil)
			Convey("Users JSON must be correctly parsed", func() {
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
				So(u.Data.LastHeartbeat.Format(time.RFC3339), ShouldEqual, "2015-04-26T04:16:23Z")
				So(u.Data.LastPlugin, ShouldEqual, "wakatime/4.0.8")
				So(u.Data.LastPluginName, ShouldEqual, "Sublime")
				So(u.Data.LastProject, ShouldEqual, "go-wakatime")
				So(u.Data.Location, ShouldEqual, "Stockholm, Sweden")
				So(u.Data.LoggedTimePublic, ShouldBeTrue)
				So(u.Data.Modified.Format(time.RFC3339), ShouldEqual, "2015-04-23T05:48:57Z")
				So(u.Data.Photo, ShouldEqual, "https://secure.gravatar.com/avatar/8cf592a20de754300721bf954aa40507?s=150&d=identicon")
				So(u.Data.PhotoPublic, ShouldBeTrue)
				So(u.Data.Plan, ShouldEqual, "basic")
				So(u.Data.Timezone, ShouldEqual, "Europe/Stockholm")
				So(u.Data.Username, ShouldEqual, "aquilax")
				So(u.Data.Website, ShouldEqual, "http://www.avtobiografia.com")
			})
		})
	})
	Convey("Given wakatime", t, func() {
		wt := New(NewDummyTransport(summaries))
		Convey("Wakatime must not be nil", func() {
			So(wt, ShouldNotBeNil)
			Convey("Summaries JSON must be correctly parsed", func() {
				s, err := wt.Summaries(CurrentUser, time.Now(), time.Now(), nil, nil)
				So(err, ShouldBeNil)
				So(s, ShouldNotBeNil)
				So(s.End.Time().Unix(), ShouldEqual, 1429912832)
				So(s.Start.Time().Unix(), ShouldEqual, 1429740032)
				So(len(s.Data), ShouldEqual, 1)
				sday := s.Data[0]

				// Editors
				So(len(sday.Editors), ShouldEqual, 1)
				So(sday.Editors[0].Digital, ShouldEqual, "2:1005")
				So(sday.Editors[0].Hours, ShouldEqual, 2)
				So(sday.Editors[0].Minutes, ShouldEqual, 10)
				So(sday.Editors[0].Name, ShouldEqual, "PhpStorm")
				So(sday.Editors[0].Percent, ShouldEqual, 69.91)
				So(sday.Editors[0].Seconds, ShouldEqual, 5)
				So(sday.Editors[0].Text, ShouldEqual, "2 hours 10 minutes 5 seconds")
				So(sday.Editors[0].TotalSeconds, ShouldEqual, 7805)

				// Grand Total
				So(sday.GrandTotal.Digital, ShouldEqual, "3:03")
				So(sday.GrandTotal.Hours, ShouldEqual, 3)
				So(sday.GrandTotal.Minutes, ShouldEqual, 3)
				So(sday.GrandTotal.Text, ShouldEqual, "3 hours 3 minutes")
				So(sday.GrandTotal.TotalSeconds, ShouldEqual, 11165)

				// Languages
				So(len(sday.Languages), ShouldEqual, 1)
				So(sday.Languages[0].Digital, ShouldEqual, "0:22:58")
				So(sday.Languages[0].Hours, ShouldEqual, 0)
				So(sday.Languages[0].Minutes, ShouldEqual, 22)
				So(sday.Languages[0].Name, ShouldEqual, "Go")
				So(sday.Languages[0].Percent, ShouldEqual, 12.34)
				So(sday.Languages[0].Seconds, ShouldEqual, 58)
				So(sday.Languages[0].Text, ShouldEqual, "22 minutes 58 seconds")
				So(sday.Languages[0].TotalSeconds, ShouldEqual, 1378)

				// Operating systems
				So(len(sday.OperatingSystems), ShouldEqual, 1)
				So(sday.OperatingSystems[0].Digital, ShouldEqual, "3:0604")
				So(sday.OperatingSystems[0].Hours, ShouldEqual, 3)
				So(sday.OperatingSystems[0].Minutes, ShouldEqual, 6)
				So(sday.OperatingSystems[0].Name, ShouldEqual, "Linux")
				So(sday.OperatingSystems[0].Percent, ShouldEqual, 100)
				So(sday.OperatingSystems[0].Seconds, ShouldEqual, 4)
				So(sday.OperatingSystems[0].Text, ShouldEqual, "3 hours 6 minutes 4 seconds")
				So(sday.OperatingSystems[0].TotalSeconds, ShouldEqual, 11164)

				// Projects
				So(len(sday.Projects), ShouldEqual, 1)
				So(sday.Projects[0].Digital, ShouldEqual, "0:22")
				So(sday.Projects[0].Hours, ShouldEqual, 0)
				So(sday.Projects[0].Minutes, ShouldEqual, 22)
				So(sday.Projects[0].Name, ShouldEqual, "go-wakatime")
				So(sday.Projects[0].Percent, ShouldEqual, 12.23)
				So(sday.Projects[0].Text, ShouldEqual, "22 minutes")
				So(sday.Projects[0].TotalSeconds, ShouldEqual, 1365)

				// Range
				So(sday.Range.Date, ShouldEqual, "04/23/2015")
				So(sday.Range.DateHuman, ShouldEqual, "04/23/2015")
				So(sday.Range.End.Time().Unix(), ShouldEqual, 1429826432)
				So(sday.Range.Start.Time().Unix(), ShouldEqual, 1429740032)
				So(sday.Range.Text, ShouldEqual, "04/23/2015")
				So(sday.Range.Timezone, ShouldEqual, "Europe/Stockholm")
			})
		})
	})

	Convey("Given wakatime", t, func() {
		wt := New(NewDummyTransport(durations))
		Convey("Wakatime must not be nil", func() {
			So(wt, ShouldNotBeNil)
			Convey("Durations JSON must be correctly parsed", func() {
				d, err := wt.Durations(CurrentUser, time.Now(), nil, nil)
				So(err, ShouldBeNil)
				So(d, ShouldNotBeNil)
				So(len(d.Branches), ShouldEqual, 1)
				So(d.Branches[0], ShouldEqual, "master")
				So(len(d.Data), ShouldEqual, 1)
				So(d.Data[0].Duration, ShouldEqual, 2240.0)
				So(d.Data[0].Project, ShouldEqual, "go-wakatime")
				So(d.Data[0].Time.Time().UnixNano(), ShouldEqual, 1430021760000000000)
				So(d.End.Time().Unix(), ShouldEqual, 1430085632)
				So(d.Start.Time().Unix(), ShouldEqual, 1429999232)
			})
		})
	})

	Convey("Given wakatime", t, func() {
		wt := New(NewDummyTransport(stats))
		Convey("Wakatime must not be nil", func() {
			So(wt, ShouldNotBeNil)
			Convey("Stats JSON must be correctly parsed", func() {
				s, err := wt.Stats(CurrentUser, Last30Days, nil, nil, nil)
				So(err, ShouldBeNil)
				So(s, ShouldNotBeNil)
				So(s.Data.CreatedAt.Format(time.RFC3339), ShouldEqual, "2015-04-23T04:38:05Z")
				// Editors
				So(len(s.Data.Editors), ShouldEqual, 1)
				So(s.Data.Editors[0].CreatedAt.Format(time.RFC3339), ShouldEqual, "2015-04-24T07:12:29Z")
				So(s.Data.Editors[0].ID, ShouldEqual, "a5979e20-aa2a-403d-9214-b49e85f00fbc")
				So(s.Data.Editors[0].ModifiedAt.Format(time.RFC3339), ShouldEqual, "2015-04-28T07:08:06Z")
				So(s.Data.Editors[0].Name, ShouldEqual, "Vim")
				So(s.Data.Editors[0].Percent, ShouldEqual, 22.64)
				So(s.Data.Editors[0].TotalSeconds, ShouldEqual, 11803)

				So(s.Data.End.Time().UnixNano(), ShouldEqual, 1430172032000000000)
				So(s.Data.HumanReadableDailyAverage, ShouldEqual, "2 hours 3 minutes")
				So(s.Data.HumanReadableTotal, ShouldEqual, "14 hours 24 minutes")
				So(s.Data.ID, ShouldEqual, "3e570b91-2540-4c9e-a71a-75b1909188ea")
				So(s.Data.IsUpToDate, ShouldBeTrue)

				// Languages
				So(len(s.Data.Languages), ShouldEqual, 1)
				So(s.Data.Languages[0].CreatedAt.Format(time.RFC3339), ShouldEqual, "2015-04-24T07:12:29Z")
				So(s.Data.Languages[0].ID, ShouldEqual, "23c0dd3f-d09a-4c01-a813-37aa257a114c")
				So(s.Data.Languages[0].ModifiedAt.Format(time.RFC3339), ShouldEqual, "2015-04-28T07:08:06Z")
				So(s.Data.Languages[0].Name, ShouldEqual, "Go")
				So(s.Data.Languages[0].Percent, ShouldEqual, 41.37)
				So(s.Data.Languages[0].TotalSeconds, ShouldEqual, 21569)

				So(s.Data.ModifiedAt.Format(time.RFC3339), ShouldEqual, "2015-04-28T07:08:06Z")

				// Operating Systems
				So(len(s.Data.OperatingSystems), ShouldEqual, 1)
				So(s.Data.OperatingSystems[0].CreatedAt.Format(time.RFC3339), ShouldEqual, "2015-04-24T07:12:29Z")
				So(s.Data.OperatingSystems[0].ID, ShouldEqual, "36685edf-5a5e-4d53-8d77-61c1dd6e9c46")
				So(s.Data.OperatingSystems[0].ModifiedAt.Format(time.RFC3339), ShouldEqual, "2015-04-28T07:08:06Z")
				So(s.Data.OperatingSystems[0].Name, ShouldEqual, "Linux")
				So(s.Data.OperatingSystems[0].Percent, ShouldEqual, 100.00)
				So(s.Data.OperatingSystems[0].TotalSeconds, ShouldEqual, 52137)

				So(s.Data.Project, ShouldBeNil)

				// Projects
				So(len(s.Data.Projects), ShouldEqual, 1)
				So(s.Data.Projects[0].CreatedAt.Format(time.RFC3339), ShouldEqual, "2015-04-24T07:12:29Z")
				So(s.Data.Projects[0].ID, ShouldEqual, "198e4ed6-a208-41b7-b698-1826a003411a")
				So(s.Data.Projects[0].ModifiedAt.Format(time.RFC3339), ShouldEqual, "2015-04-28T07:08:06Z")
				So(s.Data.Projects[0].Name, ShouldEqual, "go-wakatime")
				So(s.Data.Projects[0].Percent, ShouldEqual, 45.77)
				So(s.Data.Projects[0].TotalSeconds, ShouldEqual, 23865)

				So(s.Data.Range, ShouldEqual, Last7Days)
				So(s.Data.Start.Time().UnixNano(), ShouldEqual, 1429567232000000000)
				So(s.Data.Status, ShouldEqual, "ok")
				So(s.Data.Timeout, ShouldEqual, 15)
				So(s.Data.Timezone, ShouldEqual, "Europe/Stockholm")
				So(s.Data.TotalSeconds, ShouldEqual, 51840)
				So(s.Data.UserID, ShouldEqual, "e9b45851-991b-4755-9ccd-6355d927f472")
				So(s.Data.Username, ShouldEqual, "aquilax")
				So(s.Data.WritesOnly, ShouldBeTrue)
			})
		})
	})
	Convey("Given wakatime", t, func() {
		wt := New(NewDummyTransport(hartbeat))
		Convey("Wakatime must not be nil", func() {
			So(wt, ShouldNotBeNil)
			Convey("Durations JSON must be correctly parsed", func() {
				h, err := wt.GetHartbeats(CurrentUser, time.Now())
				So(err, ShouldBeNil)
				So(h, ShouldNotBeNil)
				So(len(h.Data), ShouldEqual, 2)
				So(h.End.Time().Format(time.RFC3339), ShouldEqual, "2015-06-03T00:00:32+02:00")
				So(h.Start.Time().Format(time.RFC3339), ShouldEqual, "2015-06-02T00:00:32+02:00")
			})
		})
	})
}
