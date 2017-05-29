package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

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
