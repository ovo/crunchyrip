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
		Usage: "download crunchyroll episodes",
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
				Usage: "comma-separated episode ID on vrv.co ex. https://vrv.co/watch/ -> GRMGEZ85R <- /Hunter-x-Hunter",
			},
			&cli.StringFlag{
				Name:  "locale",
				Value: "en-US",
			},
			&cli.StringFlag{
				Name:  "resolution",
				Value: "",
				Usage: "resolution of the download",
			},
			&cli.StringFlag{
				Name:  "seriesID",
				Value: "",
				Usage: "ID of the series you want to download a season for",
			},
		},
		Action: func(c *cli.Context) error {
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
				index, err := strconv.Atoi(strings.Trim(text, "\n"))

				if err != nil {
					return err
				}

				log.Println("Getting episode info for " + seasons[index-1].Title)
				episodes, err := cr.GetEpisodes(&client, authConfig, seasons[index-1].ID)

				if err != nil {
					return err
				}

				wg.Add(len(episodes))

				for _, e := range episodes {
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
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}

}
