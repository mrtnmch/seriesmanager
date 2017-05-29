package main

import (
	"errors"
	"strings"
)

type SeriesDetector struct {
	series map[string][][]string
}

func NewSeriesDetector(series []string) *SeriesDetector {
	detector := new(SeriesDetector)
	detector.series = make(map[string][][]string)

	for _, ser := range series {
		detector.series[ser] = make([][]string, 2)
		lower := strings.ToLower(ser)
		detector.series[ser][0] = strings.Split(lower, " ")

		temp := strings.Split(lower, " ")

		if len(temp) > 2 {
			detector.series[ser][1] = make([]string, 1)

			for i := 0; i < len(temp); i++ {
				detector.series[ser][1][0] += string(temp[i][0])
			}
		}
	}

	return detector
}

func (detector *SeriesDetector) Detect(file string) (string, error) {
	title := strings.ToLower(file)

	for ser, words := range detector.series {
	test:
		for _, set := range words {
			if len(set) == 0 {
				continue
			}

			for _, word := range set {
				if !strings.Contains(title, word) {
					continue test
				}
			}

			return ser, nil
		}

	}

	return "", errors.New("Series couln't be detected")
}