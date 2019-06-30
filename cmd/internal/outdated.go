package internal

import (
	"github.com/coreos/go-semver/semver"
)

//Outdated checks if an newer version is available
func (mod *LocalMod) Outdated() (bool, error) {
	db, err := GetGlobalMod(mod.Name)
	if err != nil {
		return false, err
	}

	newest := semver.New(db.Version)
	current := semver.New(mod.Version)
	return current.LessThan(*newest), nil
}
