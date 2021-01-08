package crunchyroll

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ovo/crunchyrip/common"
)

// Season holds information about the series season
type Season struct {
	common.Metadata
	Links struct {
		SeasonChannel struct {
			Href string `json:"href"`
		} `json:"season/channel"`
		SeasonEpisodes struct {
			Href string `json:"href"`
		} `json:"season/episodes"`
		SeasonSeries struct {
			Href string `json:"href"`
		} `json:"season/series"`
	} `json:"__links__"`
	Actions struct {
	} `json:"__actions__"`
	ID           string        `json:"id"`
	ChannelID    string        `json:"channel_id"`
	Title        string        `json:"title"`
	SeriesID     string        `json:"series_id"`
	SeasonNumber int           `json:"season_number"`
	IsComplete   bool          `json:"is_complete"`
	Description  string        `json:"description"`
	Keywords     []interface{} `json:"keywords"`
	SeasonTags   []string      `json:"season_tags"`
	Images       struct {
	} `json:"images"`
	IsMature       bool   `json:"is_mature"`
	MatureBlocked  bool   `json:"mature_blocked"`
	IsSubbed       bool   `json:"is_subbed"`
	IsDubbed       bool   `json:"is_dubbed"`
	IsSimulcast    bool   `json:"is_simulcast"`
	SeoTitle       string `json:"seo_title"`
	SeoDescription string `json:"seo_description"`
}

// GetSeasons gets the seasons information for the given seriesID
func GetSeasons(c *http.Client, auth AuthConfig, seriesID string) ([]Season, error) {
	type seasonsResp struct {
		common.Metadata
		Links struct {
		} `json:"__links__"`
		Actions struct {
		} `json:"__actions__"`
		Total int      `json:"total"`
		Items []Season `json:"items"`
	}
	var seasons seasonsResp

	req, err := http.NewRequest(http.MethodGet, "https://beta-api.crunchyroll.com/cms/v2"+auth.Bucket+"/seasons?series_id="+seriesID+"&locale=en-US&Signature="+auth.Signature+"&Key-Pair-Id="+auth.KeyPairID+"&Policy="+auth.Policy, nil)

	if err != nil {
		return []Season{}, err
	}

	req.Header.Add("User-Agent", common.UserAgent)
	req.Header.Add("Accept-Language", "en-US;q=1.0")

	resp, err := c.Do(req)

	if err != nil {
		return []Season{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []Season{}, err
	}

	json.Unmarshal([]byte(body), &seasons)

	return seasons.Items, nil
}
