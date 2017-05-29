package main

import "testing"

func TestArrayContains(t *testing.T) {
	haystack := []string{"asd", "dsa", "", " ", "...", ",", "1234654899879486541"}

	if arrayContains(haystack, "test") {
		t.Errorf("Array %s doesn't contain item 'test'", haystack)
	}

	if arrayContains(haystack, "dsa ") {
		t.Errorf("Array %s doesn't contain item 'dsa '", haystack)
	}

	if !arrayContains(haystack, "asd") {
		t.Errorf("Array %s contains item 'asd'", haystack)
	}

	if !arrayContains(haystack, "") {
		t.Errorf("Array %s contains item ''", haystack)
	}
}

func TestExtensionFilter(t *testing.T) {
	array := []string{
		"test",
		"test.mp4",
		"test.mp4.avi",
		"file.",
		"",
		".",
		"test.mp4.",
		"asd.4mp",
		"mp4.test",
		"test.mp4/file"}

	if len := len(extensionFilter(array, []string{".mp4"})); len != 1 {
		t.Errorf("Array %s contains 1 mp4 file, %d returned", array, len)
	}

	if len := len(extensionFilter(array, []string{".mp4", ".avi"})); len != 2 {
		t.Errorf("Array %s contains 2 mp4/avi files, %d returned", array, len)
	}

	if len := len(extensionFilter(array, []string{})); len != 0 {
		t.Errorf("Array should be empty, %d returned", len)
	}
}

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

func TestFilenameGenerator(t *testing.T) {
	generator := NewNameGenerator("/home/user/videos", 2)

	pairs := []struct {
		Title    string
		Season   int
		Episode  int
		Expected string
	}{
		{"The Big Bang Theory", 1, 2, "The Big Bang Theory s01e02"},
		{"How I Met Your Mother", 10, 221, "How I Met Your Mother s10e221"}}

	for _, pair := range pairs {
		if gen := generator.GenerateFilename(pair.Title, pair.Season, pair.Episode, GetExtension(pair.Expected)); gen != pair.Expected {
			t.Errorf("Name should be '%s'; '%s' returned", pair.Expected, gen)
		}
	}

	generator = NewNameGenerator("/home/user/videos", 0)
	pairs = []struct {
		Title    string
		Season   int
		Episode  int
		Expected string
	}{
		{"The Big Bang Theory", 1, 2, "The Big Bang Theory s1e2"},
		{"How I Met Your Mother", 10, 221, "How I Met Your Mother s10e221"}}

	for _, pair := range pairs {
		if gen := generator.GenerateFilename(pair.Title, pair.Season, pair.Episode, GetExtension(pair.Expected)); gen != pair.Expected {
			t.Errorf("Name should be '%s'; '%s' returned", pair.Expected, gen)
		}
	}
}

func TestFilepathGenerator(t *testing.T) {
	generator := NewNameGenerator("/home/user/videos", 2)

	pairs := []struct {
		Title    string
		Season   int
		Episode  int
		Expected string
	}{
		{"The Big Bang Theory", 1, 2, "/home/user/videos/The Big Bang Theory/Season 01/"},
		{"How I Met Your Mother", 10, 221, "/home/user/videos/How I Met Your Mother/Season 10/"}}

	for _, pair := range pairs {
		if gen := generator.GenerateFilepath(pair.Title, pair.Season, pair.Episode); gen != pair.Expected {
			t.Errorf("Name should be '%s'; '%s' returned", pair.Expected, gen)
		}
	}

	generator = NewNameGenerator("/home/user/videos", 0)
	pairs = []struct {
		Title    string
		Season   int
		Episode  int
		Expected string
	}{
		{"The Big Bang Theory", 1, 2, "/home/user/videos/The Big Bang Theory/Season 1/"},
		{"How I Met Your Mother", 10, 221, "/home/user/videos/How I Met Your Mother/Season 10/"}}

	for _, pair := range pairs {
		if gen := generator.GenerateFilepath(pair.Title, pair.Season, pair.Episode); gen != pair.Expected {
			t.Errorf("Name should be '%s'; '%s' returned", pair.Expected, gen)
		}
	}
}

func TestLocationGenerator(t *testing.T) {
	generator := NewNameGenerator("/home/user/videos", 2)

	pairs := []struct {
		Title     string
		Season    int
		Episode   int
		Extension string
		Expected  string
	}{
		{"The Big Bang Theory", 1, 2, ".mp4", "/home/user/videos/The Big Bang Theory/Season 01/The Big Bang Theory s01e02.mp4"},
		{"How I Met Your Mother", 10, 221, ".avi", "/home/user/videos/How I Met Your Mother/Season 10/How I Met Your Mother s10e221.avi"}}

	for _, pair := range pairs {
		if gen := generator.GenerateLocation(pair.Title, pair.Season, pair.Episode, pair.Extension); gen != pair.Expected {
			t.Errorf("Name should be '%s'; '%s' returned", pair.Expected, gen)
		}
	}

	generator = NewNameGenerator("/home/user/videos", 0)
	pairs = []struct {
		Title     string
		Season    int
		Episode   int
		Extension string
		Expected  string
	}{
		{"The Big Bang Theory", 1, 2, ".mp4", "/home/user/videos/The Big Bang Theory/Season 1/The Big Bang Theory s1e2.mp4"},
		{"How I Met Your Mother", 10, 221, ".avi", "/home/user/videos/How I Met Your Mother/Season 10/How I Met Your Mother s10e221.avi"}}

	for _, pair := range pairs {
		if gen := generator.GenerateLocation(pair.Title, pair.Season, pair.Episode, pair.Extension); gen != pair.Expected {
			t.Errorf("Name should be '%s'; '%s' returned", pair.Expected, gen)
		}
	}
}

func TestGetExtension(t *testing.T) {
	pairs := []struct {
		Filename string
		Expected string
	}{
		{"The Big Bang Theory.mp4", ".mp4"},
		{"/home/user/videos/The Big Bang Theory.avi", ".avi"},
		{"/this/has/no/extension", ""}}

	for _, pair := range pairs {
		if ext := GetExtension(pair.Filename); ext != pair.Expected {
			t.Errorf("Extension should be '%s', '%s' returned", pair.Expected, ext)
		}
	}
}
