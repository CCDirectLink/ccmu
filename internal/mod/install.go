package mod

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/CCDirectLink/CCUpdaterCLI/internal/moddb"
	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

//Install the mod
func (m Mod) Install() error {
	if (m.Installed()) {
		return pkg.NewErrorReason(pkg.ReasonAlreadyInstalled, pkg.ModeInstall, m, nil)
	}

	reader, size, err := moddb.DownloadMod(m.Name)
	if err != nil {
		return pkg.NewError(pkg.ModeInstall, m, err)
	}
	defer reader.Close()

	zipReader, err := zip.NewReader(newBufferedReaderAt(reader), size)
	if err != nil {
		return pkg.NewError(pkg.ModeInstall, m, err)
	}

	_, err = extract(zipReader, m.path())
	return pkg.NewError(pkg.ModeInstall, m, err)
}

func extract(reader *zip.Reader, dir string) (string, error) {
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
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
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
