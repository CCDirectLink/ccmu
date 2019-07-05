package install

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/global"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
)

//Install a mod
func Install(name string, override bool) error {
	mod, err := global.GetMod(name)
	if err != nil {
		return err
	}

	err = os.MkdirAll("installing", os.ModePerm)
	if err != nil {
		return err
	}
	defer os.RemoveAll("installing")

	file, err := download(mod.ArchiveLink)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	dir, err := extract(file)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	pkgDir, found, err := findPackage(dir)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("cmd/internal: Could not find package of mod '%s'", name)
	}

	modDir, err := getModFolderName(name, override)
	if err != nil {
		return err
	}

	if mod.Dir != nil && mod.Dir.Any == "root" {
		modDir = getRootDir(modDir)
		pkgDir = getRootDir(pkgDir)
		if !strings.HasPrefix(pkgDir, dir) {
			return fmt.Errorf("cmd/internal: Mod '%s' does not have enough directories to be installed in root", name)
		}
	}

	err = copyDir(modDir, pkgDir)
	if err != nil {
		return err
	}

	return nil
}

func findPackage(dir string) (string, bool, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return dir, false, err
	}

	for _, file := range files {
		if !file.IsDir() && file.Name() == "package.json" {
			return dir, true, nil
		}
	}

	for _, file := range files {
		if file.IsDir() {
			res, found, err := findPackage(filepath.Join(dir, file.Name()))
			if err != nil {
				return res, found, err
			}

			if found {
				return res, true, nil
			}
		}
	}

	return dir, false, nil
}

func getModFolderName(name string, override bool) (string, error) {
	path, err := local.GetGame()
	if err != nil {
		return path, err
	}

	path = filepath.Join(path, "assets", "mods", name)
	if override {
		return path, nil
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return path, nil
	}
	if err != nil {
		return path, err
	}

	for i := 2; ; i++ {
		tmpPath := path + strconv.Itoa(i)
		_, err = os.Stat(tmpPath)
		if os.IsNotExist(err) {
			return tmpPath, nil
		}
		if err != nil {
			return tmpPath, err
		}
	}
}

//getRootDir returns the third parent directory
func getRootDir(dir string) string {
	return filepath.Dir(filepath.Dir(filepath.Dir(dir)))
}
