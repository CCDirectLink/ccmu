package mod

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/CCDirectLink/ccmu/internal/moddb"
	"github.com/CCDirectLink/ccmu/pkg"
)

//Install the mod
func (m Mod) Install() error {
	if m.Installed() {
		return pkg.NewErrorReason(pkg.ReasonAlreadyInstalled, pkg.ModeInstall, m, nil)
	}

	if !m.Available() {
		return pkg.NewErrorReason(pkg.ReasonNotFound, pkg.ModeInstall, m, nil)
	}

	deps, err := m.directDeps()
	if err != nil {
		return pkg.NewError(pkg.ModeInstall, m, err)
	}

	for _, dep := range deps {
		err = dep.Install()
		var pkgErr pkg.Error
		if errors.As(err, &pkgErr) && pkgErr.Reason != pkg.ReasonAlreadyInstalled {
			return err
		} else if err != nil {
			return err
		}
	}

	buf, src, err := moddb.DownloadMod(m.Name)
	if err != nil {
		return pkg.NewError(pkg.ModeInstall, m, err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return pkg.NewError(pkg.ModeInstall, m, err)
	}

	_, err = extract(zipReader, m.path(), src)
	if err != nil {
		return pkg.NewError(pkg.ModeInstall, m, err)
	}

	return nil
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
