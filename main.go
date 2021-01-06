package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"

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
				Name:     "episodeID",
				Value:    "",
				Usage:    "episode ID on vrv.co ex. https://vrv.co/watch/ -> GRMGEZ85R <- /Hunter-x-Hunter",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "locale",
				Value: "en-US",
			},
			&cli.StringFlag{
				Name:  "resolution",
				Value: "1920x1080",
				Usage: "resolution of the download",
			},
		},
		Action: func(c *cli.Context) error {
			client := http.Client{}
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

			log.Println("Getting episode info")
			episode, err := cr.GetEpisode(&client, authConfig, c.String("episodeID"))

			if err != nil {
				return err
			}

			streamURL, err := cr.GetStreamURL(&client, authConfig, episode.Links.Streams.Href, c.String("locale"))

			if err != nil {
				return err
			}

			log.Println("Downloading stream")
			cr.DownloadStream(&client, authConfig, streamURL, c.String("resolution"))

			return nil
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}

}
