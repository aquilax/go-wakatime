package wakatime

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const ApiBase = "http://api"
const CurrentUser = "current"

type WakaTime struct {
	client *http.Client
}

type Durations struct{}

func NewWakaTime(tr *http.Transport) *WakaTime {
	return &WakaTime{
		client: &http.Client{
			Transport: tr,
		},
	}
}

func (wt *WakaTime) Durations(user string) (*Durations, error) {
	var err error
	var resp *http.Response
	if resp, err = wt.fetch("users/" + user + "/durations"); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var content []byte
	if content, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	var dr Durations
	if err = json.Unmarshal(content, dr); err != nil {
		return nil, err
	}
	return &dr, nil
}

func (wt *WakaTime) Stats() {}

func (wt *WakaTime) Summaries() {}

func (wt *WakaTime) Users() {}

func (wt *WakaTime) fetch(path string) (*http.Response, error) {
	return wt.client.Get(wt.getURL(path))
}

func (wt *WakaTime) getURL(path string) string {
	return ApiBase + path
}
