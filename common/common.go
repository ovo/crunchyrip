package common

import "strings"

const UserAgent = "Crunchyroll/4.3.0 (bundle_identifier:com.crunchyroll.iphone; build_number:1832556.306266127) iOS/14.2.0 Gravity/3.0.0"

// FormatTitle formats the titles of episodes, seasons, and series
func FormatTitle(s string) string {
	return strings.Join(strings.Split(s, " "), "_")
}
