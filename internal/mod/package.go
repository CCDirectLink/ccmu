package mod

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/CCDirectLink/ccmu/internal/mod/installer"
)

//Package represents the data inside package.json
type Package struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Version      string            `json:"version"`
	Licence      string            `json:"licence"`
	Dependencies map[string]string `json:"ccmodDependencies"`
	Hidden       bool              `json:"hidden"`
}

func readPackage(reader io.Reader) (Package, error) {
	var result Package
	err := json.NewDecoder(reader).Decode(&result)
	return result, err
}

func (m Mod) readPackageFile() (Package, error) {
	if filepath.Ext(m.path()) == ".ccmod" {
		return readPackageFromCCMod(m.path())
	}

	path := filepath.Join(m.path(), "package.json")
	file, err := os.Open(path)
	if err != nil {
		return Package{}, err
	}
	defer file.Close()

	return readPackage(file)
}

func readPackageFromSource(source string) (Package, error) {
	tmp := filepath.Join(os.TempDir(), "ccmu")
	os.MkdirAll(tmp, os.ModeDir)
	defer os.RemoveAll(tmp)

	result := Package{
		Name:    source,
		Version: "0.0.0",
	}

	err := installer.Packed{
		ModDir: tmp,
		Path:   source,
	}.Install()
	if err != nil {
		return result, err
	}

	files, err := ioutil.ReadDir(tmp)
	if err != nil {
		return result, err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".ccmod" {
			return readPackageFromCCMod(filepath.Join(tmp, file.Name()))
		}
	}

	return result, os.ErrNotExist
}

func readPackageFromCCMod(path string) (Package, error) {
	file, err := os.Open(path)
	if err != nil {
		return Package{}, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return Package{}, err
	}

	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return Package{}, err
	}

	return readPackedPackage(reader)
}

func readPackedPackage(reader *zip.Reader) (Package, error) {
	for _, file := range reader.File {
		path := filepath.Join("/", file.Name)

		if len(path) > 0 && path[1:] == "package.json" && !file.FileInfo().IsDir() {
			stream, err := file.Open()
			defer stream.Close()
			if err != nil {
				return Package{}, err
			}
			return readPackage(stream)
		}
	}
	return Package{}, os.ErrNotExist
}

func extract(reader *zip.Reader, dst, src string) (string, error) {
	for _, file := range reader.File {
		if len(file.Name) < len(src) || file.Name[:len(src)] != src {
			continue
		}

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dst, file.Name[len(src):])

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dst)+string(os.PathSeparator)) && (fpath != filepath.Clean(dst) || !file.FileInfo().IsDir()) {
			return dst, fmt.Errorf("%s: illegal file path", fpath)
		}

		if file.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return dst, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return dst, err
		}

		tmpFile, err := file.Open()
		if err != nil {
			return dst, err
		}

		_, err = io.Copy(outFile, tmpFile)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		tmpFile.Close()

		if err != nil {
			return dst, err
		}
	}
	return dst, nil
}
