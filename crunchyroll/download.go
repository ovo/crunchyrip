package crunchyroll

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/grafov/m3u8"
	lru "github.com/hashicorp/golang-lru"
	"github.com/ovo/crunchyrip/common"
)

var (
	// IVPlaceholder holds place for IV
	IVPlaceholder = []byte{0, 0, 0, 0, 0, 0, 0, 0}
)

// Download holds information needed for downloading
type Download struct {
	URI           string
	SeqNo         uint64
	ExtXKey       *m3u8.Key
	totalDuration time.Duration
}

// DownloadStream downloads the given stream url
func DownloadStream(c *http.Client, auth AuthConfig, url string, resolution string, ep Episode, wg *sync.WaitGroup) error {
	defer wg.Done()
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", common.UserAgent)

	resp, err := c.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	os.Mkdir("./downloads", 0755)
	os.Mkdir("./downloads/"+common.FormatTitle(ep.SeriesTitle), 0755)
	os.Mkdir("./downloads/"+common.FormatTitle(ep.SeriesTitle)+"/"+common.FormatTitle(ep.SeasonTitle), 0755)
	filePath := "./downloads/" + common.FormatTitle(ep.SeriesTitle) + "/" + common.FormatTitle(ep.SeasonTitle) + "/" + common.FormatTitle(ep.Title) + "_" + ep.ID + ".ts"

	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), true)

	if err != nil {
		return err
	}

	if listType == m3u8.MASTER {
		masterpl := p.(*m3u8.MasterPlaylist)

		if resolution != "" {
			for i, variant := range masterpl.Variants {
				if variant.Resolution == resolution {
					msChan := make(chan *Download, 1024)
					go GetPlaylist(c, masterpl.Variants[i].URI, 0, true, msChan, wg)
					DownloadSegment(c, filePath, msChan, 0)
					break
				}
			}
		} else {
			msChan := make(chan *Download, 1024)
			go GetPlaylist(c, masterpl.Variants[0].URI, 0, true, msChan, wg)
			DownloadSegment(c, filePath, msChan, 0)
		}

	}
	return nil
}

// DecryptData decrypts the AES-128 encrypted data
func DecryptData(c *http.Client, data []byte, v *Download, aes128Keys *map[string][]byte) {
	var (
		iv          *bytes.Buffer
		keyData     []byte
		cipherBlock cipher.Block
	)

	if v.ExtXKey != nil && (v.ExtXKey.Method == "AES-128" || v.ExtXKey.Method == "aes-128") {

		keyData = (*aes128Keys)[v.ExtXKey.URI]

		if keyData == nil {
			req, _ := http.NewRequest("GET", v.ExtXKey.URI, nil)
			resp, _ := c.Do(req)
			keyData, _ = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			(*aes128Keys)[v.ExtXKey.URI] = keyData
		}

		if v.ExtXKey.IV == "" {
			iv = bytes.NewBuffer(IVPlaceholder)
			binary.Write(iv, binary.BigEndian, v.SeqNo)
		} else {
			iv = bytes.NewBufferString(v.ExtXKey.IV)
		}

		cipherBlock, _ = aes.NewCipher((*aes128Keys)[v.ExtXKey.URI])
		cipher.NewCBCDecrypter(cipherBlock, iv.Bytes()).CryptBlocks(data, data)
	}

}

// DownloadSegment downloads the segment of the file
func DownloadSegment(c *http.Client, fn string, dlc chan *Download, recTime time.Duration) {
	var out, err = os.Create(fn)
	defer out.Close()

	if err != nil {
		log.Fatal(err)
		return
	}
	var (
		data       []byte
		aes128Keys = &map[string][]byte{}
	)

	for v := range dlc {
		req, err := http.NewRequest("GET", v.URI, nil)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := c.Do(req)
		if err != nil {
			log.Print(err)
			continue
		}
		if resp.StatusCode != 200 {
			log.Printf("Received HTTP %v for %v\n", resp.StatusCode, v.URI)
			continue
		}

		data, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		DecryptData(c, data, v, aes128Keys)

		_, err = out.Write(data)

		// _, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Downloaded %v\n", v.URI)
		if recTime != 0 {
			log.Printf("Recorded %v of %v\n", v.totalDuration, recTime)
		} else {
			log.Printf("Recorded %v\n", v.totalDuration)
		}
	}
}

// GetPlaylist gets the playlist data
func GetPlaylist(c *http.Client, urlStr string, recTime time.Duration, useLocalTime bool, dlc chan *Download, wg *sync.WaitGroup) {
	startTime := time.Now()
	var recDuration time.Duration = 0
	cache, _ := lru.New(1024)
	playlistURL, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := c.Do(req)
		if err != nil {
			log.Print(err)
			time.Sleep(time.Duration(3) * time.Second)
		}
		playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		if listType == m3u8.MEDIA {
			mpl := playlist.(*m3u8.MediaPlaylist)

			for segmentIndex, v := range mpl.Segments {
				if v != nil {
					var msURI string
					if strings.HasPrefix(v.URI, "http") {
						msURI, err = url.QueryUnescape(v.URI)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						msURL, err := playlistURL.Parse(v.URI)
						if err != nil {
							log.Print(err)
							continue
						}
						msURI, err = url.QueryUnescape(msURL.String())
						if err != nil {
							log.Fatal(err)
						}
					}
					_, hit := cache.Get(msURI)
					if !hit {
						cache.Add(msURI, nil)
						if useLocalTime {
							recDuration = time.Now().Sub(startTime)
						} else {
							recDuration += time.Duration(int64(v.Duration * 1000000000))
						}
						dlc <- &Download{
							URI:           msURI,
							ExtXKey:       mpl.Key,
							SeqNo:         uint64(segmentIndex) + mpl.SeqNo,
							totalDuration: recDuration,
						}
					}
					if recTime != 0 && recDuration != 0 && recDuration >= recTime {
						close(dlc)
						return
					}
				}
			}
			if mpl.Closed {
				close(dlc)
				return
			}
			time.Sleep(time.Duration(int64(mpl.TargetDuration * 1000000000)))

		} else {
			log.Fatal("Not a valid media playlist")
		}
	}
}
