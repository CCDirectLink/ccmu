package installer

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
