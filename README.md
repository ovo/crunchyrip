[![Go Report Card](https://goreportcard.com/badge/github.com/ovo/crunchyrip)](https://goreportcard.com/report/github.com/ovo/crunchyrip) [![CircleCI](https://circleci.com/gh/ovo/crunchyrip.svg?style=svg)](https://circleci.com/gh/ovo/crunchyrip) ![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
# Crunchyrip

Download full episodes from crunchyroll into a .ts media file

Inspired by [anirip](https://github.com/s32x/anirip)

## Dependencies
- Go
- ffmpeg

## Installation
Clone or download repository

`$ go install`

## Usage

#### For individual episodes

`$ crunchyrip download --email <email> --password <password> --episodeIDs <episodeID>,<episodeID2>...`

#### For seasons

`$ crunchyrip download --email <email> --password <password> --seriesID <seriesID>`

You will be prompted to select the season you want to download for the given series

Episodes will be stored in the downloads folder

For more info run `$ crunchyrip [subcommand] --help`

#### Episode range

`$ crunchyrip download --email <email> --password <password> --seriesID <seriesID> --range <start episodeID>-<end episodeID>`

This is useful for when you want to download multiple episodes but do not want to comma-seperate every episode or download the entire season

## Find episodeID and seriesID

**If you are on beta crunchyroll, the ID should be in the url of the episode or season**

## Finding episodeID

1. Go to Crunchyroll and find the episode you want to download
2. Inspect element and paste this into console
`document.getElementsByClassName('boxcontents')[0].id.split('_')[2]`

## Finding seriesID

1. Go to Crunchyroll and find the series that you want to download
2. Inspect element and paste this into console
`JSON.parse(document.getElementsByClassName("show-actions")[0].attributes['data-contentmedia'].value).mediaId`
