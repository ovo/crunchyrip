[![Go Report Card](https://goreportcard.com/badge/github.com/ovo/crunchyrip)](https://goreportcard.com/report/github.com/ovo/crunchyrip) ![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
# Crunchyrip

Download full episodes from crunchyroll into a .ts media file

Inspired by [anirip](https://github.com/s32x/anirip)

## Installation
Clone or download repository

`$ go install`

## Usage

`$ crunchyrip --email <email> --password <password> --episodeID <episodeID>,<episodeID2>...`

Episodes will be stored in the downloads folder

For more info run `$ crunchyrip --help`

## Finding episodeID

1. Go to https://vrv.co and find the episode you want to download
2. For the episode https://vrv.co/watch/GRMGEZ85R/Hunter-x-Hunter, the ID would be GRMGEZ85R