package wakatime

type WakaTime struct {
	APIKey string
}

func NewWakaTime(APIKey string) *WakaTime {
	return &WakaTime{APIKey}
}

func (wt *WakaTime) Durations() {}
func (wt *WakaTime) Stats()     {}
func (wt *WakaTime) Summaries() {}
func (wt *WakaTime) Users()     {}
