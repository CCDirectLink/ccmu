package local

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//GetGame using the current working directory or flags
func GetGame() (string, error) {
	dir, err := getDir()
	if err != nil {
		return "", err
	}
	return searchForGame(dir)
}

func getDir() (string, error) {
	game := flag.Lookup("game")
	if game != nil {
		return game.Value.String(), nil
	}

	return os.Getwd()
}

func searchForGame(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	if containsPackage(files) {
		if exists, _ := exists(filepath.Join(dir, "./assets/node-webkit.html")); exists {
			return dir, nil
		}
	}

	parent := filepath.Dir(dir)
	if parent == dir {
		return "", fmt.Errorf("cmd/internal: Could not find game")
	}

	return searchForGame(parent)
}

func containsPackage(files []os.FileInfo) bool {
	for _, file := range files {
		if !file.IsDir() && file.Name() == "package.json" {
			return true
		}
	}
	return false
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if !os.IsNotExist(err) {
		return true, err
	}
	return false, nil
}
