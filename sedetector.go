package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type SeasonEpisodeDetector struct {
	regexes []*regexp.Regexp
}

func NewSeasonEpisodeDetector() *SeasonEpisodeDetector {
	detector := new(SeasonEpisodeDetector)

	detector.regexes = []*regexp.Regexp{
		regexp.MustCompile("(.*)s(?P<season>\\d+)(\\s*)e(?P<episode>\\d+)([^\\d]*)"),
		regexp.MustCompile("([^\\d]*)(?P<season>\\d+)x(?P<episode>\\d+)([^\\d]*)")}

	return detector
}

func (detector *SeasonEpisodeDetector) Detect(file string) (season, episode int, err error) {
	for _, reg := range detector.regexes {
		match := reg.FindStringSubmatch(strings.ToLower(file))

		if match == nil {
			continue
		}

		result := make(map[string]string)
		for i, name := range reg.SubexpNames() {
			if i != 0 {
				result[name] = match[i]
			}
		}

		season, e := strconv.Atoi(result["season"])
		episode, e2 := strconv.Atoi(result["episode"])

		if e != nil {
			return 0, 0, e
		}

		if e2 != nil {
			return 0, 0, e2
		}

		return season, episode, nil
	}

	return 0, 0, errors.New("Season/episode pattern not detected")
}