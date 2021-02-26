package common

import "strings"

// UserAgent is the user agent used by the crunchyroll app
const UserAgent = "Crunchyroll/4.3.0 (bundle_identifier:com.crunchyroll.iphone; build_number:1832556.306266127) iOS/14.5.0 Gravity/3.0.0"

// ClientID is used for authentication - hard coded into the app
const ClientID = "lrses_25zh3kjta6r9r9"

// ClientSecret is used for authentication - hard coded into the app
const ClientSecret = "afDXiDZ8a6ElzKQpF6o2psLbcG3iZFHs"

// FormatTitle formats the titles of episodes, seasons, and series
func FormatTitle(s string) string {
	return strings.Join(strings.Split(s, " "), "_")
}

// Metadata hold metadata information about the response
type Metadata struct {
	Class       string `json:"__class__"`
	Href        string `json:"__href__"`
	ResourceKey string `json:"__resource_key__"`
}
