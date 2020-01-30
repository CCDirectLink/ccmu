package tool

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/CCDirectLink/ccmu/internal/game"
	"github.com/CCDirectLink/ccmu/internal/mod"
	"github.com/CCDirectLink/ccmu/internal/moddb"
	"github.com/CCDirectLink/ccmu/pkg"
)

type ccloader struct {
	game game.Game
}

func (c ccloader) Info() (pkg.Info, error) {
	info, errOnline := moddb.PkgInfo("ccloader")
	path, err := c.game.BasePath()
	if err != nil {
		return info, nil
	}

	ppath := filepath.Join(path, "ccloader", "package.json")
	existing, err := exists(ppath)
	if err != nil {
		return info, pkg.NewError(pkg.ModeUnknown, c, err)
	}

	if existing {
		err = c.readPackage(ppath, &info)
	} else {
		err = c.readJS(path, &info)
	}
	if err != nil {
		return info, pkg.NewError(pkg.ModeUnknown, c, err)
	}

	if errOnline != nil {
		return info, pkg.NewError(pkg.ModeUnknown, c, errOnline)
	}
	return info, nil
}

func (c ccloader) readPackage(path string, info *pkg.Info) error {
	var data moddb.PackageDBPackageMetadata

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&data)
	if info.Name == "" {
		info.Name = data.Name
	}
	if info.NiceName == "" {
		info.NiceName = data.CCModHumanName
	}
	if info.Description == "" {
		info.Description = data.Description
	}
	if info.Licence == "" {
		info.Licence = data.Licence
	}
	info.CurrentVersion = string(data.Version)

	return err
}

func (c ccloader) readJS(base string, info *pkg.Info) error {
	if info.Name == "" {
		info.Name = "ccloader"
	}
	if info.NiceName == "" {
		info.NiceName = "CCLoader"
	}
	if info.Description == "" {
		info.Description = "Modloader for CrossCode. This or a similar modloader is needed for most mods."
	}

	path := filepath.Join(base, "ccloader", "js", "ccloader")
	existing, err := exists(path)
	if err != nil {
		return err
	}
	if !existing {
		return pkg.ErrNotFound
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	defer reader.Reset(nil)
	reader.ReadSlice('=')
	reader.ReadSlice('\'')
	res, err := reader.ReadString('\'')
	if err != nil {
		return err
	}

	if len(res) > 0 {
		res = res[:len(res)-1]
	}
	info.CurrentVersion = res

	return nil
}

func (c ccloader) Installed() bool {
	base, err := c.game.BasePath()
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

	base, err := c.game.BasePath()
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
	base, err := c.game.BasePath()
	if err != nil {
		return pkg.NewError(pkg.ModeUninstall, c, err)
	}

	var result error
	err = os.RemoveAll(filepath.Join(base, "ccloader"))
	if err != nil {
		result = pkg.NewError(pkg.ModeUninstall, c, err)
	}

	file, err := os.Create(filepath.Join(base, "package.json"))
	if err != nil {
		result = pkg.NewError(pkg.ModeUninstall, c, err)
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

	return result
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
	base, err := c.game.BasePath()
	result := []pkg.Package{
		mod.Mod{
			Name: "Simplify",
			Path: base,
			Game: c.game,
		},
	}

	if err != nil {
		return result, pkg.NewError(pkg.ModeUnknown, c, err)
	}
	return result, nil
}

func (c ccloader) NewestDependencies() ([]pkg.Package, error) {
	return c.Dependencies()
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
