package cmd

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/global"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
	"github.com/Masterminds/semver"
)

func installDependencies(mod local.Mod, stats *Stats) error {
	for name, version := range mod.Dependencies {
		if err := installDependency(name, version, stats); err != nil {
			return err
		}
	}
	return nil
}

func installDependency(name, version string, stats *Stats) error {
	ver, err := semver.NewConstraint(version)
	if err != nil {
		stats.AddWarning(fmt.Sprintf("cmd: Mod '%s' had an invalid version number '%s'", name, version))
		return nil
	}

	newest, err := global.GetMod(name)
	if err != nil {
		stats.AddWarning(fmt.Sprintf("cmd: Could find mod '%s'", name))
		return nil
	}

	newestVer, err := semver.NewVersion(newest.Version)
	if err != nil {
		stats.AddWarning(fmt.Sprintf("cmd: Could not parse mod list: Mod '%s' has an invalid version number '%s'", name, version))
		return nil
	}

	if !ver.Check(newestVer) {
		stats.AddWarning(fmt.Sprintf("cmd: Could not update mod '%s' to %s because the newest version is %s", name, version, newest.Version))
		return nil
	}

	mod, err := local.GetMod(name)
	if err != nil {
		return installMod(name, stats)
	}

	outdated, err := mod.Outdated()
	if err != nil {
		return err
	}

	if outdated {
		return updateMod(name, stats)
	}
	return nil
}
