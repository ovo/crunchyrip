package crunchyroll

import (
	"io/ioutil"
	"net/http"

	"github.com/Jeffail/gabs/v2"
	"github.com/ovo/crunchyrip/common"
)

// GetStreamURL returns the stream url for the videoID and locale
func GetStreamURL(c *http.Client, auth AuthConfig, path string, locale string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://beta-api.crunchyroll.com"+path+"?locale="+locale+"&Key-Pair-Id="+auth.KeyPairID+"&Policy="+auth.Policy+"&Signature="+auth.Signature, nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", common.UserAgent)
	req.Header.Add("Accept-Language", "en-US;q=1.0")

	resp, err := c.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	var url string

	body, _ := ioutil.ReadAll(resp.Body)

	jsonParsed, err := gabs.ParseJSON([]byte(body))

	if err != nil {
		return "", err
	}

	url, _ = jsonParsed.Path("streams.adaptive_hls." + locale + ".url").Data().(string)

	return url, nil
}
