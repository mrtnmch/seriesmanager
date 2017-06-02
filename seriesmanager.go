package main

import (
	"flag"
	"fmt"
	"path"
)

const CONFIG_FILE = ".sm-config.json"

func main() {
	var createSymlinks, testRun, silent, force, moved, renameOnly bool

	flag.BoolVar(&createSymlinks, "symlink", false, "Create a symlink in the original location")
	flag.BoolVar(&testRun, "test", false, "Do only a test run, don't move files or create symlinks")
	flag.BoolVar(&silent, "silent", false, "Silent run - do not print any output")
	flag.BoolVar(&force, "force", false, "Force move - even if the target file exists (override)")
	flag.BoolVar(&renameOnly, "rename", false, "Rename only - do not move")
	flag.Parse()

	config, err := LoadConfig(relativePath(CONFIG_FILE))

	if err != nil {
		fmt.Println(err)
		config = CreateDefaultConfig()
		config.Save(relativePath(CONFIG_FILE))
	}

	for _, inputPath := range config.InputPaths {
		files, err := listFiles(inputPath, true)

		if err != nil {
			fmt.Printf("Error while loading '%s'\n", inputPath)
			continue
		}

		files = extensionFilter(files, config.Extensions)

		detector := config.seriesDetector
		sedetector := config.sessionEpisodeDetector
		generator := config.nameGenerator

		for _, file := range files {
			det, err := detector.Detect(file)
			season, episode, err2 := sedetector.Detect(file)

			if err == nil && err2 == nil {
				if renameOnly {
					generator = NewNameGenerator(path.Dir(file), generator.pad, true)
				}

				gen := generator.GenerateLocation(det, season, episode, getExtension(file))

				if gen == file {
					continue
				}

				exists := exists(gen)

				if exists && !force {
					if !silent {
						fmt.Printf("Skipping %s\n%s exists\n---\n", file, gen)
					}

					continue
				}

				moved = true

				if !testRun {
					makePath(generator.GenerateFilepath(det, season, episode))
					move(file, gen)

					if createSymlinks {
						createSymlink(gen, file)
					}
				}

				if !silent {
					fmt.Printf("%s\n%s\n\n", file, gen)
				}
			}
		}
	}

	if testRun && !silent {
		fmt.Println("Test run - no files moved or renamed")
	}

	if !moved && !silent {
		fmt.Println("No files moved or renamed, everything's where it belongs")
	}
}
