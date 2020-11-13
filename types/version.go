package types

import "time"

var appVersion = "?"
var buildTime = "?"
var commit = "?"

// Utility structure
type Version struct {
	App         string `json:"app"`         // Application.
	AppVersion  string `json:"app_version"` // Application version.
	BuildTime   string `json:"build_time"`  // Time when application was build.
	Commit      string `json:"commit"`      // Git commit hash.
	TimeZone    string `json:"time_zone"`
	CurrentTime string `json:"current_time"`
}

func NewVersion(app string) *Version {
	return &Version{
		App:         app,
		AppVersion:  appVersion,
		BuildTime:   buildTime,
		Commit:      commit,
		TimeZone:    time.Local.String(),
		CurrentTime: time.Now().Round(0).String(),
	}
}
