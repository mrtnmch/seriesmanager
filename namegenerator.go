package main

import (
	"fmt"
	"path"
	"strconv"
	"strings"
)

const titleWildcard = "<title>"
const seasonWildcard = "<season>"
const episodeWildcard = "<episode>"

type NameGenerator struct {
	pattern     string
	pad         int
	pathPattern string
	basePath    string
}

func NewNameGenerator(basePath string, pad int, renameOnly bool) *NameGenerator {
	generator := new(NameGenerator)
	generator.pattern = fmt.Sprintf("%s s%se%s", titleWildcard, seasonWildcard, episodeWildcard)
	generator.basePath = basePath
	generator.pad = pad

	if renameOnly {
		generator.pathPattern = "/"
	} else {
		generator.pathPattern = fmt.Sprintf("/%s/S%s/", titleWildcard, seasonWildcard)
	}

	return generator
}

func leftPad(s string, padStr string, pLen int) string {
	if pLen-len(s) < 0 {
		return s
	}

	return strings.Repeat(padStr, pLen-len(s)) + s
}

func (generator *NameGenerator) GenerateFilename(title string, season, episode int, extension string) string {
	ret := generator.pattern
	ret = strings.Replace(ret, titleWildcard, title, -1)
	ret = strings.Replace(ret, seasonWildcard, leftPad(strconv.Itoa(season), "0", generator.pad), -1)
	ret = strings.Replace(ret, episodeWildcard, leftPad(strconv.Itoa(episode), "0", generator.pad), -1)
	ret = ret + extension
	return ret
}

func (generator *NameGenerator) GenerateFilepath(title string, season, episode int) string {
	ret := generator.basePath + generator.pathPattern
	ret = strings.Replace(ret, titleWildcard, title, -1)
	ret = strings.Replace(ret, seasonWildcard, leftPad(strconv.Itoa(season), "0", generator.pad), -1)
	ret = strings.Replace(ret, episodeWildcard, leftPad(strconv.Itoa(episode), "0", generator.pad), -1)
	return ret
}

func (generator *NameGenerator) GenerateLocation(title string, season, episode int, extension string) string {
	return path.Join(generator.GenerateFilepath(title, season, episode), generator.GenerateFilename(title, season, episode, extension))
}
