package main

import (
	"encoding/json"
	"io/ioutil"
	"os/user"
	"path"
)

type Config struct {
	seriesDetector         *SeriesDetector
	sessionEpisodeDetector *SeasonEpisodeDetector
	nameGenerator          *NameGenerator
	InputPaths             []string
	OutputPath             string
	Series                 []string
	PadNumbers             int
	Extensions             []string
}

func CreateDefaultConfig() *Config {
	config := new(Config)

	myself, err := user.Current()
	if err != nil {
		return nil
	}

	config.InputPaths = []string{myself.HomeDir}
	config.OutputPath = path.Join(myself.HomeDir, "Videos")
	config.Series = []string{"The Big Bang Theory", "How I Met Your Mother", "11.22.63", "13 Reasons Why", "Billions", "Dexter", "Friends", "Game of Thrones", "Homeland", "House of Cards", "Mr. Robot", "Narcos", "Prison Break", "Shooter", "Suits", "Stargate Universe", "The Grand Tour", "The Man In The High Castle", "Westworld", "Silicon Valley"}
	config.PadNumbers = 2
	config.Extensions = []string{".avi", ".mp4", ".mkv", ".wmv", ".srt", ".sub"}

	config.seriesDetector = NewSeriesDetector(config.Series)
	config.sessionEpisodeDetector = NewSeasonEpisodeDetector()
	config.nameGenerator = NewNameGenerator(config.OutputPath, config.PadNumbers, false)

	return config
}

func (config *Config) Save(path string) error {
	out, err := json.Marshal(config)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, out, 0644)

	return err
}

func (config *Config) Reload() {
	config.seriesDetector = NewSeriesDetector(config.Series)
	config.sessionEpisodeDetector = NewSeasonEpisodeDetector()
	config.nameGenerator = NewNameGenerator(config.OutputPath, config.PadNumbers, false)
}

func LoadConfig(path string) (*Config, error) {
	config := CreateDefaultConfig()
	out, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, config)

	if err != nil {
		return nil, err
	}

	config.Reload()

	return config, err
}
