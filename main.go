package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"sync"

	cr "github.com/ovo/crunchyrip/crunchyroll"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:  "crunchyrip",
		Usage: "download full crunchyroll episodes",
		Commands: []*cli.Command{
			{
				Name:    "download",
				Aliases: []string{"d"},
				Usage:   "download episodes or a season",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "email",
						Value:    "",
						Usage:    "email for crunchyroll account",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password",
						Value:    "",
						Usage:    "password for the crunchyroll account",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "episodeIDs",
						Value: "",
						Usage: "comma-separated episode IDs; no comma needed for single episode",
					},
					&cli.StringFlag{
						Name:  "locale",
						Value: "en-US",
					},
					&cli.StringFlag{
						Name:  "resolution",
						Value: "",
						Usage: "resolution of the download (ex. 1080x1920) - defaults to highest resolution",
					},
					&cli.StringFlag{
						Name:  "seriesID",
						Value: "",
						Usage: "ID of the series you want to download a season for",
					},
					&cli.StringFlag{
						Name:  "range",
						Value: "",
						Usage: "combine with seriesID to download a range of episodes (ex. --range G69P9MD9Y-GRGGQ42DR)",
					},
				},
				Action: downloadAction,
			},
			{
				Name:    "resolution",
				Aliases: []string{"r"},
				Usage:   "get resolutions for an episode",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "email",
						Value:    "",
						Usage:    "email for crunchyroll account",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password",
						Value:    "",
						Usage:    "password for the crunchyroll account",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "episodeID",
						Value:    "",
						Usage:    "episodeID to get resolutions for",
						Required: true,
					},
				},
				Action: resolutionAction,
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}

}

func downloadAction(c *cli.Context) error {
	var wg sync.WaitGroup
	tr := &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 20,
	}
	client := http.Client{Transport: tr}
	client.Jar, _ = cookiejar.New(nil)

	log.Println("Logging in")
	credentials, err := cr.Login(&client, c.String("email"), c.String("password"))

	if err != nil {
		return err
	}

	log.Println("Getting CMS info")
	cms, err := cr.GetCMS(&client, credentials.AccessToken)

	if err != nil {
		return err
	}

	authConfig := cr.AuthConfig{
		AccessToken: credentials.AccessToken,
		Policy:      cms.Cms.Policy,
		Signature:   cms.Cms.Signature,
		KeyPairID:   cms.Cms.KeyPairID,
		Bucket:      cms.Cms.Bucket,
	}

	if c.String("episodeIDs") != "" {
		ids := strings.Split(c.String("episodeIDs"), ",")

		wg.Add(len(ids))

		for _, id := range ids {
			go func(c *http.Client, authConfig cr.AuthConfig, id string, resolution string, locale string) {
				log.Println("Getting episode info for " + id)

				episode, err := cr.GetEpisode(c, authConfig, id)

				if err != nil {
					log.Fatal(err)
				}

				streamURL, err := cr.GetStreamURL(c, authConfig, episode.Links.Streams.Href, locale)

				if err != nil {
					log.Fatal(err)
				}
				log.Println("Downloading " + episode.ID)
				go cr.DownloadStream(c, authConfig, streamURL, resolution, episode, &wg)

			}(&client, authConfig, id, c.String("resolution"), c.String("locale"))

		}
	}

	if c.String("seriesID") != "" {
		seasons, err := cr.GetSeasons(&client, authConfig, c.String("seriesID"))

		if err != nil {
			return err
		}

		for i, s := range seasons {
			fmt.Println(i+1, s.Title)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter season to download: ")
		text, _ := reader.ReadString('\n')
		index, err := strconv.Atoi(strings.TrimSpace(text))

		if err != nil {
			return err
		}

		log.Println("Getting episode info for " + seasons[index-1].Title)
		episodes, err := cr.GetEpisodes(&client, authConfig, seasons[index-1].ID)

		if err != nil {
			return err
		}

		var newep []cr.Episode

		if c.String("range") != "" {
			found := false
			eprange := strings.Split(c.String("range"), "-")
			for _, e := range episodes {
				if e.ID == eprange[0] {
					found = true
				}
				if found {
					newep = append(newep, e)
				}
				if e.ID == eprange[1] {
					found = false
				}
			}
		} else {
			newep = episodes
		}

		for _, e := range newep {
			wg.Add(1)
			go func(c *http.Client, authConfig cr.AuthConfig, ep cr.Episode, resolution string, locale string) {
				streamURL, err := cr.GetStreamURL(c, authConfig, ep.Links.Streams.Href, locale)

				if err != nil {
					log.Fatal(err)
				}

				log.Println("Downloading " + ep.ID)
				go cr.DownloadStream(c, authConfig, streamURL, resolution, ep, &wg)
			}(&client, authConfig, e, c.String("resolution"), c.String("locale"))
		}
	}

	wg.Wait()

	return nil
}

func resolutionAction(c *cli.Context) error {
	tr := &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 20,
	}
	id := c.String("episodeID")
	client := http.Client{Transport: tr}
	client.Jar, _ = cookiejar.New(nil)

	log.Println("Logging in")
	credentials, err := cr.Login(&client, c.String("email"), c.String("password"))

	if err != nil {
		return err
	}

	log.Println("Getting CMS info")
	cms, err := cr.GetCMS(&client, credentials.AccessToken)

	authConfig := cr.AuthConfig{
		AccessToken: credentials.AccessToken,
		Policy:      cms.Cms.Policy,
		Signature:   cms.Cms.Signature,
		KeyPairID:   cms.Cms.KeyPairID,
		Bucket:      cms.Cms.Bucket,
	}

	log.Println("Getting episode info for " + id)

	episode, err := cr.GetEpisode(&client, authConfig, id)

	if err != nil {
		return err
	}

	streamURL, err := cr.GetStreamURL(&client, authConfig, episode.Links.Streams.Href, c.String("locale"))

	if err != nil {
		return err
	}

	resolutions, err := cr.GetResolutions(&client, authConfig, streamURL, episode)

	fmt.Println()
	for _, resolution := range resolutions {
		fmt.Println(resolution)
	}

	return nil
}
