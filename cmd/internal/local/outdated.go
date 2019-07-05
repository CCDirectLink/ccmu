package local

import (
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/global"
	"github.com/coreos/go-semver/semver"
)

//Outdated checks if an newer version is available
func (mod *Mod) Outdated() (bool, error) {
	db, err := global.GetMod(mod.Name)
	if err != nil {
		return false, err
	}

	newest := semver.New(db.Version)
	current := semver.New(mod.Version)
	return current.LessThan(*newest), nil
}
