package tool

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	ccmu "github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/internal/mod"
	"github.com/CCDirectLink/CCUpdaterCLI/internal/moddb"
	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

type ccloader struct {
	path string
}

func (c ccloader) Info() (pkg.Info, error) {
	return moddb.PkgInfo("ccloader") //TODO: add local info
}

func (c ccloader) Installed() bool {
	base, err := c.base()
	if err != nil {
		return false
	}

	stat, err := os.Stat(filepath.Join(base, "ccloader"))
	return err == nil && stat.IsDir()
}

func (c ccloader) Available() bool {
	_, err := moddb.PkgInfo("ccloader")
	return err == nil
}

func (c ccloader) Install() error {
	if c.Installed() {
		return pkg.NewErrorReason(pkg.ReasonAlreadyInstalled, pkg.ModeInstall, c, nil)
	}

	if !c.Available() {
		return pkg.NewErrorReason(pkg.ReasonNotFound, pkg.ModeInstall, c, nil)
	}

	base, err := c.base()
	if err != nil {
		return pkg.NewError(pkg.ModeInstall, c, err)
	}

	buf, src, err := moddb.DownloadMod("ccloader")
	if err != nil {
		return pkg.NewError(pkg.ModeInstall, c, err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return pkg.NewError(pkg.ModeInstall, c, err)
	}

	_, err = extract(zipReader, base, src)
	if err != nil {
		return pkg.NewError(pkg.ModeInstall, c, err)
	}

	return nil
}

func (c ccloader) Uninstall() error {
	base, err := c.base()
	if err != nil {
		return pkg.NewError(pkg.ModeUninstall, c, err)
	}

	err = os.RemoveAll(filepath.Join(base, "ccloader"))
	if err != nil {
		return pkg.NewError(pkg.ModeUninstall, c, err)
	}

	file, err := os.Open(filepath.Join(base, "package.json"))
	if err != nil {
		return pkg.NewError(pkg.ModeUninstall, c, err)
	}
	defer file.Close()

	_, err = file.Write([]byte(`
	{
		"name": "CrossCode",
		"version" : "1.0.0",
		"main": "assets/node-webkit.html",
		"chromium-args" : "--ignore-gpu-blacklist --disable-direct-composition --disable-background-networking --in-process-gpu --password-store=basic",
		"window" : {
			"toolbar" : false,
			"icon" : "favicon.png",
			"width" : 1136,
			"height": 640,
			"fullscreen" : false
		}
	}`))
	if err != nil {
		return pkg.NewError(pkg.ModeUninstall, c, err)
	}

	return nil
}

func (c ccloader) Update() error {
	err := c.Uninstall()
	if err != nil {
		if pkgErr, ok := err.(pkg.Error); ok {
			pkgErr.Mode = pkg.ModeUpdate
			return pkgErr
		}
		return err
	}
	err = c.Install()
	if err != nil {
		if pkgErr, ok := err.(pkg.Error); ok {
			pkgErr.Mode = pkg.ModeUpdate
			return pkgErr
		}
	}
	return err
}

func (c ccloader) Dependencies() ([]pkg.Package, error) {
	base, err := c.base()
	result := []pkg.Package{
		mod.Mod{
			Name: "Simplify",
			Path: base,
			Game: ccmu.At(base),
		},
	}

	if err != nil {
		return result, pkg.NewError(pkg.ModeUnknown, c, err)
	}
	return result, nil
}

func (c ccloader) base() (string, error) {
	return c.path, nil //TODO: Implement
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
