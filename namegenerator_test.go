package main

import "testing"

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
		{"The Big Bang Theory", 1, 2, "/home/user/videos/The Big Bang Theory/S01/"},
		{"How I Met Your Mother", 10, 221, "/home/user/videos/How I Met Your Mother/S10/"}}

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
		{"The Big Bang Theory", 1, 2, "/home/user/videos/The Big Bang Theory/S1/"},
		{"How I Met Your Mother", 10, 221, "/home/user/videos/How I Met Your Mother/S10/"}}

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
		{"The Big Bang Theory", 1, 2, ".mp4", "/home/user/videos/The Big Bang Theory/S01/The Big Bang Theory s01e02.mp4"},
		{"How I Met Your Mother", 10, 221, ".avi", "/home/user/videos/How I Met Your Mother/S10/How I Met Your Mother s10e221.avi"}}

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
		{"The Big Bang Theory", 1, 2, ".mp4", "/home/user/videos/The Big Bang Theory/S1/The Big Bang Theory s1e2.mp4"},
		{"How I Met Your Mother", 10, 221, ".avi", "/home/user/videos/How I Met Your Mother/S10/How I Met Your Mother s10e221.avi"}}

	for _, pair := range pairs {
		if gen := generator.GenerateLocation(pair.Title, pair.Season, pair.Episode, pair.Extension); gen != pair.Expected {
			t.Errorf("Name should be '%s'; '%s' returned", pair.Expected, gen)
		}
	}
}