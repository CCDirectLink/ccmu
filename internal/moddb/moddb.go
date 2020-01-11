package moddb

import (
	"bytes"
	"io"
	"net/http"

	"github.com/CCDirectLink/ccmu/pkg"
)

var (
	pkgs           PackageDB
	pkgsDownloaded bool

	infos       map[string]pkg.Info
	infosParsed bool
)

//PkgInfo reads a given pkg.Info from the moddb.
func PkgInfo(name string) (pkg.Info, error) {
	if !infosParsed {
		_, err := PkgInfos()
		if err != nil {
			return pkg.Info{Name: name}, err
		}
	}

	result, ok := infos[name]
	if !ok {
		return result, pkg.ErrNotFound
	}
	return result, nil
}

//PkgInfos reads all pkg.Infos from the moddb.
func PkgInfos() ([]pkg.Info, error) {
	if infosParsed {
		result := make([]pkg.Info, len(infos))
		i := 0
		for _, info := range infos {
			result[i] = info
			i++
		}
		return result, nil
	}

	db, err := packageDB()
	if err != nil {
		return nil, err
	}

	cache := make(map[string]pkg.Info)
	result := make([]pkg.Info, len(db))
	i := 0
	for name, p := range db {
		result[i] = pkg.Info{
			Name:          p.Metadata.Name,
			NiceName:      p.Metadata.niceName(),
			Description:   p.Metadata.Description,
			Licence:       p.Metadata.Licence,
			NewestVersion: string(p.Metadata.Version),
			Hidden:        p.Metadata.Hidden,
		}
		cache[name] = result[i]
		i++
	}

	infos = cache
	infosParsed = true

	return result, nil
}

//MergePkgInfo with old one.
func MergePkgInfo(info *pkg.Info) error {
	newInfo, err := PkgInfo(info.Name)
	if err != nil {
		return err
	}

	info.NiceName = newInfo.NiceName
	info.Description = newInfo.Description
	info.Licence = newInfo.Licence
	info.NewestVersion = newInfo.NewestVersion
	info.Hidden = newInfo.Hidden

	return nil
}

//DownloadMod as io.ReadCloser.
func DownloadMod(name string) (bytes.Buffer, string, error) {
	p, err := packageByName(name)
	if err != nil {
		return bytes.Buffer{}, "", err
	}

	//TODO: iterate over installation method
	data, err := p.Installation[0].modZip()
	if err != nil {
		return bytes.Buffer{}, "", err
	}

	resp, err := http.Get(data.URL)
	if err != nil {
		return bytes.Buffer{}, "", err
	}
	defer resp.Body.Close()

	var result bytes.Buffer
	_, err = io.Copy(&result, resp.Body)
	return result, data.Source, err
}

//Dependencies of a package.
func Dependencies(name string) (map[string]string, error) {
	pkg, err := packageByName(name)
	result := make(map[string]string, len(pkg.Metadata.CCModDependencies))
	for k, v := range pkg.Metadata.CCModDependencies {
		result[k] = string(v)
	}
	return result, err
}

func packageDB() (PackageDB, error) {
	if pkgsDownloaded {
		return pkgs, nil
	}

	var err error
	pkgs, err = getPackageDB()
	if err != nil {
		return pkgs, err
	}

	pkgsDownloaded = true
	return pkgs, nil
}

func packageByName(name string) (PackageDBPackage, error) {
	pkgDB, err := packageDB()
	if err != nil {
		return PackageDBPackage{}, err
	}

	p, ok := pkgDB[name]
	if !ok {
		return p, pkg.ErrNotFound
	}
	return p, nil
}
