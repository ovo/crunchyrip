package crunchyroll

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Episode contains stream information about the episode
type Episode struct {
	Class       string `json:"__class__"`
	Href        string `json:"__href__"`
	ResourceKey string `json:"__resource_key__"`
	Links       struct {
		EpisodeChannel struct {
			Href string `json:"href"`
		} `json:"episode/channel"`
		EpisodeNextEpisode struct {
			Href string `json:"href"`
		} `json:"episode/next_episode"`
		EpisodeSeason struct {
			Href string `json:"href"`
		} `json:"episode/season"`
		EpisodeSeries struct {
			Href string `json:"href"`
		} `json:"episode/series"`
		Streams struct {
			Href string `json:"href"`
		} `json:"streams"`
	} `json:"__links__"`
	Actions struct {
	} `json:"__actions__"`
	ID                  string   `json:"id"`
	ChannelID           string   `json:"channel_id"`
	SeriesID            string   `json:"series_id"`
	SeriesTitle         string   `json:"series_title"`
	SeasonID            string   `json:"season_id"`
	SeasonTitle         string   `json:"season_title"`
	SeasonNumber        int      `json:"season_number"`
	Episode             string   `json:"episode"`
	EpisodeNumber       int      `json:"episode_number"`
	SequenceNumber      int      `json:"sequence_number"`
	ProductionEpisodeID string   `json:"production_episode_id"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	NextEpisodeID       string   `json:"next_episode_id"`
	NextEpisodeTitle    string   `json:"next_episode_title"`
	HdFlag              bool     `json:"hd_flag"`
	IsMature            bool     `json:"is_mature"`
	MatureBlocked       bool     `json:"mature_blocked"`
	EpisodeAirDate      string   `json:"episode_air_date"`
	IsSubbed            bool     `json:"is_subbed"`
	IsDubbed            bool     `json:"is_dubbed"`
	IsClip              bool     `json:"is_clip"`
	SeoTitle            string   `json:"seo_title"`
	SeoDescription      string   `json:"seo_description"`
	SeasonTags          []string `json:"season_tags"`
	AvailableOffline    bool     `json:"available_offline"`
	MediaType           string   `json:"media_type"`
	Slug                string   `json:"slug"`
	Images              struct {
		Thumbnail [][]struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Type   string `json:"type"`
			Source string `json:"source"`
		} `json:"thumbnail"`
	} `json:"images"`
	DurationMs      int      `json:"duration_ms"`
	IsPremiumOnly   bool     `json:"is_premium_only"`
	ListingID       string   `json:"listing_id"`
	SubtitleLocales []string `json:"subtitle_locales"`
	Playback        string   `json:"playback"`
}

// GetEpisode returns information about the episode
func GetEpisode(c *http.Client, auth AuthConfig, videoID string) (Episode, error) {
	url := "https://beta-api.crunchyroll.com/cms/v2" + auth.Bucket + "/episodes/" + videoID + "?locale=en-US&Signature=" + auth.Signature + "&Key-Pair-Id=" + auth.KeyPairID + "&Policy=" + auth.Policy
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return Episode{}, err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept-Language", "en-US;q=1.0")

	resp, err := c.Do(req)

	if err != nil {
		log.Fatal(err)
		return Episode{}, err
	}

	defer resp.Body.Close()

	var episode Episode
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return Episode{}, err
	}

	json.Unmarshal([]byte(body), &episode)

	return episode, nil
}
