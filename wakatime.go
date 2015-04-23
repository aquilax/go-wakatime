package wakatime

import (
	"net/http"
)

const ApiBase = "http://api"

type WakaTime struct {
	client *http.Client
}

func NewWakaTime(tr *http.Transport) *WakaTime {
	return &WakaTime{
		client: &http.Client{
			Transport: tr,
		},
	}
}

func (wt *WakaTime) Durations() {
	_, _ = wt.fetch("/durations")
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
