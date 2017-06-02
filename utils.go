package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/kardianos/osext"
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

func getExtension(file string) string {
	return path.Ext(file)
}

func createSymlink(file, link string) error {
	return os.Symlink(file, link)
}

func makePath(path string) error {
	return os.MkdirAll(path, 0766)
}

func move(from, to string) error {
	return os.Rename(from, to)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getExecPath() string {
	ex, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Path: %s\n", ex)

	return ex
}

func relativePath(file string) string {
	return path.Join(getExecPath(), file)
}
