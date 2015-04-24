package wakatime

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const Version = "0.1"
const ApiBase = "https://wakatime.com/api/v1/"
const CurrentUser = "current"

type WakaTime struct {
	client *http.Client
}

type Durations struct{}

func NewWakaTime(rt http.RoundTripper) *WakaTime {
	return &WakaTime{
		client: &http.Client{
			Transport: rt,
		},
	}
}

func (wt *WakaTime) Durations(user string, date time.Time, project, branches *string) (*Durations, error) {
	var err error
	var resp *http.Response
	var u *url.URL
	if u, err = url.Parse(ApiBase); err != nil {
		return nil, err
	}
	u.Path += "users/" + user + "/durations"
	q := u.Query()
	q.Set("date", date.Format("01/02/2006"))
	if project != nil {
		q.Set("project", *project)
	}
	if branches != nil {
		q.Set("branches", *branches)
	}
	u.RawQuery = q.Encode()
	if resp, err = wt.client.Get(u.String()); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var content []byte
	if content, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	fmt.Println(string(content))
	var dr Durations
	if err = json.Unmarshal(content, dr); err != nil {
		return nil, err
	}
	return &dr, nil
}

func (wt *WakaTime) Stats() {}

func (wt *WakaTime) Summaries() {}

func (wt *WakaTime) Users() {}

func (wt *WakaTime) getURL(path string) string {
	return ApiBase + path
}
