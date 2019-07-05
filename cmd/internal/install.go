package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//Install a mod
func Install(name string, override bool) error {
	mod, err := GetGlobalMod(name)
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

func download(url string) (*os.File, error) {
	file, err := ioutil.TempFile("installing", "mod")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return file, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return file, err
}

func extract(file *os.File) (string, error) {
	dir, err := ioutil.TempDir("installing", "mod")
	if err != nil {
		return "", err
	}

	reader, err := zip.OpenReader(file.Name())
	if err != nil {
		return "", err
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dir, file.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dir)+string(os.PathSeparator)) {
			return dir, fmt.Errorf("%s: illegal file path", fpath)
		}

		if file.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return dir, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return dir, err
		}

		tmpFile, err := file.Open()
		if err != nil {
			return dir, err
		}

		_, err = io.Copy(outFile, tmpFile)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		tmpFile.Close()

		if err != nil {
			return dir, err
		}
	}
	return dir, nil
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
	path, err := GetGame()
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
