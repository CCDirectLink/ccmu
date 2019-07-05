package install

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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
