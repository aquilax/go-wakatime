package wakatime

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Version is the library version
const Version = "0.1"

// APIBase is the root URL of the API
const APIBase = "https://wakatime.com/api/v1/"

// CurrentUser replaces the current user in the request
const CurrentUser = "current"

const dateFormat = "01/02/2006"

// Range is the stats report interval range
type Range string

// Stats report ranges
const (
	Last7Days   Range = "last_7_days"
	Last30Days        = "last_30_days"
	Last6Months       = "last_6_months"
	LastYear          = "last_year"
	AllTime           = "all_time"
)

// Time is time.Time alias, used for parsing response timestamps
type Time time.Time

// WakaTime is the main structure
type WakaTime struct {
	client *http.Client
}

// DurationsData is single duration segment
type DurationsData struct {
	Duration Time
	Project  string
	Time     Time
}

// Durations is the structure rerurned by the durations request
type Durations struct {
	Branches []string
	Data     []DurationsData
	End      Time
	Start    Time
	TimeZone string
}

// StatsItem is single item in the stats report
type StatsItem struct {
	CreatedAt    time.Time
	ID           string
	ModifiedAt   time.Time
	Name         string
	Percent      float32
	TotalSeconds int
}

// StatsEditor represents editor data in the stats report
type StatsEditor StatsItem

// StatsLanguage represents language data in the stats report
type StatsLanguage StatsItem

// StatsOperatingSystem represents operating system data in the stats report
type StatsOperatingSystem StatsItem

// StatsProject represents project data in the stats report
type StatsProject StatsItem

// StatsData is the main data body in the stats report
type StatsData struct {
	CreatedAt                 time.Time
	Editors                   []StatsEditor
	End                       Time
	HumanReadableDailyAverage string
	HumanReadableTotal        string
	ID                        string
	IsUpToDate                bool
	Languages                 []StatsLanguage
	ModifiedAt                time.Time
	OperatingSystems          []StatsOperatingSystem
	Project                   *string
	Projects                  []StatsProject
	Range                     Range
	Start                     Time
	Status                    string
	Timeout                   int
	Timezone                  string
	TotalSeconds              int
	UserID                    string
	Username                  string
	WritesOnly                bool
}

// Stats is the data returned by the stats report
type Stats struct {
	Data StatsData
}

// SummaryGrandTotal contains the grand total summary
type SummaryGrandTotal struct {
	Digital      string
	Hours        int
	Minutes      int
	Seconds      int
	Text         string
	TotalSeconds int
}

// SummaryItem contains the summary item data
type SummaryItem struct {
	SummaryGrandTotal
	Name    string
	Percent float32
}

type SummaryEditor SummaryItem
type SummaryLanguage SummaryItem
type SummaryOperatingSystem SummaryItem
type SummaryProject SummaryItem

type SumaryRange struct {
	Date      string
	DateHuman string
	End       Time
	Start     Time
	Text      string
	Timezone  string
}

type SummariesData struct {
	Editors          []SummaryEditor
	GrandTotal       SummaryGrandTotal
	Languages        []SummaryLanguage
	OperatingSystems []SummaryOperatingSystem
	Projects         []SummaryProject
	Range            SumaryRange
}

type Summaries struct {
	Data  []SummariesData
	End   Time
	Start Time
}

type UserData struct {
	Created              time.Time
	Email                string
	EmailPublic          bool   `json:"email_public"`
	FullName             string `json:"full_name"`
	HumanReadableWebsite string `json:"human_readable_website"`
	ID                   string
	LastHeartbeat        time.Time `json:"last_heartbeat"`
	LastPlugin           string    `json:"last_plugin"`
	LastPluginName       string    `json:"last_plugin_name"`
	LastProject          string    `json:"last_project"`
	Location             string
	LoggedTimePublic     bool `json:"logged_time_public"`
	Modified             time.Time
	Photo                string
	PhotoPublic          bool `json:"photo_public"`
	Plan                 string
	Timezone             string
	Username             string
	Website              string
}

type Users struct {
	Data UserData
}

// New initializes the library
func New(rt http.RoundTripper) *WakaTime {
	return &WakaTime{
		client: &http.Client{
			Transport: rt,
		},
	}
}

// Durations fetches the durations report
func (wt *WakaTime) Durations(user string, date time.Time, project, branches *string) (*Durations, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "users/" + user + "/durations"
	q := u.Query()
	q.Set("date", date.Format(dateFormat))
	if project != nil {
		q.Set("project", *project)
	}
	if branches != nil {
		q.Set("branches", *branches)
	}
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = wt.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var dr Durations
	if err = json.Unmarshal(content, &dr); err != nil {
		return nil, err
	}
	return &dr, nil
}

// Stats fetches the stats report
func (wt *WakaTime) Stats(user string, rng Range, timeout *int, writesOnly *bool, project *string) (*Stats, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "users/" + user + "/stats/" + rng.String()
	q := u.Query()
	if timeout != nil {
		q.Set("timeout", strconv.Itoa(*timeout))
	}
	if writesOnly != nil {
		q.Set("writes_only", strconv.FormatBool(*writesOnly))
	}
	if project != nil {
		q.Set("project", *project)
	}
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = wt.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var st Stats
	if err = json.Unmarshal(content, &st); err != nil {
		return nil, err
	}
	return &st, nil
}

// Summaries fetches the summaries report
func (wt *WakaTime) Summaries(user string, start, end time.Time, project, branches *string) (*Summaries, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "users/" + user + "/summaries"
	q := u.Query()
	q.Set("start", start.Format(dateFormat))
	q.Set("end", start.Format(dateFormat))
	if project != nil {
		q.Set("project", *project)
	}
	if branches != nil {
		q.Set("branches", *branches)
	}
	u.RawQuery = q.Encode()
	var content []byte
	if content, err = wt.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var sm Summaries
	if err = json.Unmarshal(content, &sm); err != nil {
		return nil, err
	}
	return &sm, nil
}

// Users fetches the users report
func (wt *WakaTime) Users(user string) (*Users, error) {
	var err error
	var u *url.URL
	if u, err = url.Parse(APIBase); err != nil {
		return nil, err
	}
	u.Path += "users/" + user
	var content []byte
	if content, err = wt.fetchURL(u.String()); err != nil {
		return nil, err
	}
	var us Users
	if err = json.Unmarshal(content, &us); err != nil {
		return nil, err
	}
	return &us, nil
}

// UnmarshalJSON unmarshals the Time type
func (ut *Time) UnmarshalJSON(data []byte) error {
	ts, err := strconv.ParseFloat(string(data), 32)
	if err != nil {
		return err
	}
	sec := int64(ts)
	ns := int64((ts - float64(sec)) * float64(time.Second))
	*ut = Time(time.Unix(int64(sec), ns))
	return nil
}

// String returns the string representation of Range
func (r Range) String() string {
	return string(r)
}

func (wt *WakaTime) fetchURL(url string) ([]byte, error) {
	var err error
	var resp *http.Response
	if resp, err = wt.client.Get(url); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Error: %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
