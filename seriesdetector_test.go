package main

import "testing"

func TestSeriesDetector(t *testing.T) {
	detector := NewSeriesDetector([]string{"The Big Bang Theory", "How I Met Your Mother"})
	pairs := []struct {
		Test     string
		Expected string
	}{
		{"the big bang theory s04e01.mp4", "The Big Bang Theory"},
		{"tbbt s04e01.mp4", "The Big Bang Theory"},
		{"HOW I met YOUR MOther asd.avi", "How I Met Your Mother"}}

	for _, pair := range pairs {
		if det, _ := detector.Detect(pair.Test); det != pair.Expected {
			t.Errorf("'%s' should be '%s'; '%s' returned", pair.Test, pair.Expected, det)
		}
	}
}