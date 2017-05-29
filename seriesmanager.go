package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const titleWildcard = "<title>"
const seasonWildcard = "<season>"
const episodeWildcard = "<episode>"

func listFiles(path string, recursive bool) ([]string, error) {
	files, _ := ioutil.ReadDir(path)

	var ret []string

	for _, file := range files {
		abs, err := filepath.Abs(filepath.Join(path, file.Name()))

		if err != nil {
			return nil, err
		}

		if !file.IsDir() {
			ret = append(ret, abs)
		} else if recursive {
			temp, err := listFiles(abs, recursive)

			if err != nil {
				return nil, err
			}

			for _, tFile := range temp {
				ret = append(ret, tFile)
			}
		}
	}

	return ret, nil
}

func arrayContains(array []string, needle string) bool {
	for _, cmp := range array {
		if cmp == needle {
			return true
		}
	}

	return false
}

func extensionFilter(files []string, extensions []string) []string {
	var ret []string

	for _, file := range files {
		if arrayContains(extensions, filepath.Ext(file)) {
			ret = append(ret, file)
		}
	}

	return ret
}

func print(files []string) {
	for _, file := range files {
		fmt.Println(file)
	}
}

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

type NameGenerator struct {
	pattern     string
	pad         int
	pathPattern string
	basePath    string
}

func NewNameGenerator(basePath string, pad int) *NameGenerator {
	generator := new(NameGenerator)
	generator.pattern = fmt.Sprintf("%s s%se%s", titleWildcard, seasonWildcard, episodeWildcard)
	generator.pad = pad
	generator.pathPattern = fmt.Sprintf("/%s/S%s/", titleWildcard, seasonWildcard)
	generator.basePath = basePath
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

func GetExtension(file string) string {
	return path.Ext(file)
}

func CreateSymlink(file, link string) error {
	return os.Symlink(file, link)
}

func MakePath(path string) error {
	return os.MkdirAll(path, 0766)
}

func Move(from, to string) error {
	return os.Rename(from, to)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func main() {
	var createSymlinks, testRun, silent, force, moved bool

	flag.BoolVar(&createSymlinks, "symlink", false, "Create a symlink in the original location")
	flag.BoolVar(&testRun, "test", false, "Do only a test run, don't move files or create symlinks")
	flag.BoolVar(&silent, "silent", false, "Silent run - do not print any output")
	flag.BoolVar(&force, "force", false, "Force move - even if the target file exists (override)")
	flag.Parse()

	files, _ := listFiles("/home/mxmx/downloads", true)
	files = extensionFilter(files, []string{".avi", ".mp4", ".mkv", ".wmv", ".srt", ".sub"})

	if len(files) == 0 {
		if !silent {
			fmt.Println("No matching files found, nothing to do")
		}

		return
	}

	detector := NewSeriesDetector([]string{"The Big Bang Theory", "How I Met Your Mother", "11.22.63", "13 Reasons Why", "Billions", "Dexter", "Friends", "Game of Thrones", "Homeland", "House of Cards", "Mr. Robot", "Narcos", "Prison Break", "Shooter", "Suits", "Stargate Universe", "The Grand Tour", "The Man In The High Castle", "Westworld", "Silicon Valley"})
	sedetector := NewSeasonEpisodeDetector()
	generator := NewNameGenerator("/home/mxmx/data/Videa/Original/Seriály/", 2)

	for _, file := range files {
		det, err := detector.Detect(file)
		season, episode, err2 := sedetector.Detect(file)

		if err == nil && err2 == nil {
			gen := generator.GenerateLocation(det, season, episode, GetExtension(file))

			if gen == file {
				continue
			}

			exists := Exists(gen)

			if exists && !force {
				if !silent {
					fmt.Printf("Skipping %s\n%s exists\n---\n", file, gen)
				}

				continue
			}

			moved = true

			if !testRun {
				MakePath(generator.GenerateFilepath(det, season, episode))
				Move(file, gen)

				if createSymlinks {
					CreateSymlink(gen, file)
				}
			}

			if !silent {
				fmt.Printf("%s\n%s\n---\n", file, gen)
			}
		}
	}

	if !moved && !silent {
		fmt.Println("No files moved, everything's where it belongs")
	}

	return
}
