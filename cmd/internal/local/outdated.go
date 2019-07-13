package local

import (
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/global"
	"github.com/Masterminds/semver"
)

//Outdated checks if an newer version is available
func (mod *Mod) Outdated() (bool, error) {
	db, err := global.GetMod(mod.Name)
	if err != nil {
		return false, err
	}

	newest, err := semver.NewVersion(db.Version)
	if err != nil {
		return false, err
	}
	current, err := semver.NewVersion(mod.Version)
	if err != nil {
		return false, err
	}

	return current.LessThan(newest), nil
}
