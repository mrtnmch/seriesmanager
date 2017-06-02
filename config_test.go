package main

import "testing"
import "reflect"
import "os"

func TestCreateDefaultConfig(t *testing.T) {
	config := CreateDefaultConfig()

	if config == nil {
		t.Errorf("Default config should be created, nil returned")
	}
}

func TestSaveLoadConfig(t *testing.T) {
	config := CreateDefaultConfig()
	file := ".sm-test.json"
	if err := config.Save(file); err != nil {
		t.Errorf("Couldn't save the config file")
	}

	config2, err := LoadConfig(file)

	if err != nil {
		t.Errorf("Couldn't load the saved config file")
	}

	if !reflect.DeepEqual(config, config2) {
		t.Errorf("Saved and loaded configs should be the same, a different config loaded")
	}

	os.Remove(file)
}
