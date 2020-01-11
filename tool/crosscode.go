package tool

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/CCDirectLink/ccmu/internal/game"
	"github.com/CCDirectLink/ccmu/pkg"
)

type crosscode struct {
	game game.Game
}

func (c crosscode) Info() (pkg.Info, error) {
	version, err := c.readVersion()
	result := pkg.Info{
		Name:           "crosscode",
		NiceName:       "CrossCode",
		Description:    "The base version of CrossCode",
		CurrentVersion: version,
		NewestVersion:  version,
		Hidden:         false,
	}

	if err != nil {
		return result, pkg.NewError(pkg.ModeUnknown, c, err)
	}
	return result, nil
}

func (c crosscode) Installed() bool {
	_, err := c.game.BasePath()
	return err == nil
}

func (c crosscode) Available() bool {
	return false
}

func (c crosscode) Install() error {
	return pkg.NewErrorReason(pkg.ReasonNotAvailable, pkg.ModeInstall, c, nil)
}

func (c crosscode) Uninstall() error {
	return pkg.NewErrorReason(pkg.ReasonNotAvailable, pkg.ModeInstall, c, nil)
}

func (c crosscode) Update() error {
	return pkg.NewErrorReason(pkg.ReasonNotAvailable, pkg.ModeInstall, c, nil)
}

func (c crosscode) Dependencies() ([]pkg.Package, error) {
	return []pkg.Package{}, nil
}

func (c crosscode) NewestDependencies() ([]pkg.Package, error) {
	return []pkg.Package{}, nil
}

func (c crosscode) readVersion() (string, error) {
	path, err := c.game.BasePath()
	if err != nil {
		return "0.0.0", err
	}

	file, err := os.Open(filepath.Join(path, "assets", "data", "changelog.json"))
	if err != nil {
		return "0.0.0", err
	}
	defer file.Close()

	var data struct {
		Changelog []struct {
			Version string `json:"version"`
		} `json:"changelog"`
	}

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return "0.0.0", err
	}

	if len(data.Changelog) == 0 {
		return "0.0.0", pkg.ErrNotFound
	}
	return data.Changelog[0].Version, nil
}
