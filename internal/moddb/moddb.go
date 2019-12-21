package moddb

import (
	"errors"
	"io"
	"net/http"

	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

var (
	mods           Mods
	modsDownloaded bool

	//ErrNotFound is returned when the given mod is not found.
	ErrNotFound = errors.New("moddb: Mod not found")
)

//ModInfo reads a given pkg.Info from the moddb.
func ModInfo(name string) (pkg.Info, error) {
	data, err := modInfo(name)
	if err != nil {
		return pkg.Info{}, err
	}
	return pkg.Info{
		Name:          name,
		NiceName:      data.Name,
		Description:   data.Description,
		Licence:       data.Licence,
		NewestVersion: data.Version,
	}, nil
}

//ModInfos reads all pkg.Infos from the moddb.
func ModInfos() ([]pkg.Info, error) {
	data, err := modInfos()
	if err != nil {
		return nil, err
	}

	result := make([]pkg.Info, len(data.Mods))
	i := 0
	for name, mod := range data.Mods {
		result[i] = pkg.Info{
			Name:          name,
			NiceName:      mod.Name,
			Description:   mod.Description,
			Licence:       mod.Licence,
			NewestVersion: mod.Version,
		}
		i++
	}
	return result, nil
}

//MergeModInfo with old one.
func MergeModInfo(info *pkg.Info) error {
	newInfo, err := ModInfo(info.Name)
	if err != nil {
		return err
	}

	info.NiceName = newInfo.NiceName
	info.Description = newInfo.Description
	info.Licence = newInfo.Licence
	info.NewestVersion = newInfo.NewestVersion

	return nil
}

//DownloadMod as io.ReadCloser.
func DownloadMod(name string) (io.ReadCloser, int64, error) {
	data, err := modInfo(name)
	if err != nil {
		return nil, 0, err
	}

	resp, err := http.Get(data.ArchiveLink)
	return resp.Body, resp.ContentLength, err
}

func modInfos() (Mods, error) {
	if modsDownloaded {
		return mods, nil
	}

	var err error
	mods, err = getMods()
	if err != nil {
		return mods, err
	}

	modsDownloaded = true
	return mods, nil
}

func modInfo(name string) (Mod, error) {
	mods, err := modInfos()
	if err != nil {
		return Mod{}, err
	}

	mod, ok := mods.Mods[name]
	if !ok {
		return mod, ErrNotFound
	}
	return mod, nil
}
