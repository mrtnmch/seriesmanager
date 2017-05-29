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
