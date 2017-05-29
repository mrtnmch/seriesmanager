package main

import "testing"

func TestSeasonEpisodeDetector(t *testing.T) {
	detector := NewSeasonEpisodeDetector()

	pairs := []struct {
		Test            string
		ExpectedSeason  int
		ExpectedEpisode int
	}{
		{"the big bang theory s04e01.mp4", 4, 1},
		{"the big bang theory S04E01.mp4", 4, 1},
		{"tbbt 12x10.mp4", 12, 10},
		{"HOW I met YOUR MOther asd.avi", 0, 0}}

	for _, pair := range pairs {
		if season, episode, _ := detector.Detect(pair.Test); season != pair.ExpectedSeason || episode != pair.ExpectedEpisode {
			t.Errorf("'%s' should be '%d:%d'; '%d:%d' returned", pair.Test, pair.ExpectedSeason, pair.ExpectedEpisode, season, episode)
		}
	}
}