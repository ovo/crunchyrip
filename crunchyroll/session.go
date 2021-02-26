package crunchyroll

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ovo/crunchyrip/common"
)

// Credentials holds crunchyroll oauth credientials
type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	Country      string `json:"country"`
}

// CMSInfo holds CMS info for crunchyroll streaming
type CMSInfo struct {
	Cms struct {
		Bucket    string    `json:"bucket"`
		Policy    string    `json:"policy"`
		Signature string    `json:"signature"`
		KeyPairID string    `json:"key_pair_id"`
		Expires   time.Time `json:"expires"`
	} `json:"cms"`
	ServiceAvailable bool `json:"service_available"`
}

// AuthConfig stores information needed for authorized requests
type AuthConfig struct {
	AccessToken string
	Policy      string
	Signature   string
	KeyPairID   string
	Bucket      string
}

// Login creates a new crunchyroll session
func Login(c *http.Client, user string, pass string) (Credentials, error) {
	data := url.Values{
		"grant_type": {"password"},
		"password":   {pass},
		"scope":      {"account content offline_access"},
		"username":   {user},
	}
	reader := strings.NewReader(data.Encode())
	req, err := http.NewRequest(http.MethodPost, "https://beta-api.crunchyroll.com/auth/v1/token", reader)
	authString := base64.StdEncoding.EncodeToString([]byte(common.ClientID + ":" + common.ClientSecret))

	if err != nil {
		return Credentials{}, err
	}

	req.Header.Add("User-Agent", common.UserAgent)
	req.Header.Add("Authorization", "Basic "+authString)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Add("Accept-Language", "en-US;q=1.0")

	resp, err := c.Do(req)

	if err != nil {
		return Credentials{}, err
	}

	defer resp.Body.Close()

	var credientials Credentials
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &credientials)

	return credientials, nil
}

// GetCMS gets crunchyroll CMS info
func GetCMS(c *http.Client, token string) (CMSInfo, error) {
	req, err := http.NewRequest(http.MethodGet, "https://beta-api.crunchyroll.com/index/v2", nil)

	if err != nil {
		return CMSInfo{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", common.UserAgent)
	req.Header.Add("Accept-Language", "en-US;q=1.0")

	resp, err := c.Do(req)

	if err != nil {
		return CMSInfo{}, err
	}

	defer resp.Body.Close()

	var cms CMSInfo
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &cms)

	return cms, nil
}
