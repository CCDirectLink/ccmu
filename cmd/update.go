package cmd

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/global"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/install"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
)

//Update a mod
func Update(args []string) (*Stats, error) {
	if _, err := local.GetGame(); err != nil {
		return nil, fmt.Errorf("cmd: Could not find game folder")
	}

	_, err := global.FetchModData()
	if err != nil {
		return nil, fmt.Errorf("cmd: Could not download mod data because an error occured in %s", err.Error())
	}

	if len(args) == 0 {
		return updateOutdated()
	}

	stats := &Stats{}
	for _, name := range args {
		if err := updateMod(name, stats); err != nil {
			return stats, err
		}
	}

	return stats, nil
}

func updateOutdated() (*Stats, error) {
	mods, err := local.GetMods()
	if err != nil {
		return nil, fmt.Errorf("cmd: Could not list installed mods because and error occured in %s", err.Error())
	}

	stats := &Stats{}
	for _, mod := range mods {
		if _, err := global.GetMod(mod.Name); err != nil {
			continue
		}

		outdated, err := mod.Outdated()
		if err != nil {
			return stats, fmt.Errorf("cmd: Could not check if the mod was outdated because an error occured in %s", err.Error())
		}

		if !outdated {
			continue
		}

		if err := updateMod(mod.Name, stats); err != nil {
			return stats, err
		}
	}

	return stats, nil
}

func updateMod(name string, stats *Stats) error {
	if _, err := local.GetMod(name); err != nil {
		stats.AddWarning(fmt.Sprintf("cmd: Could not update '%s' because it was not installed", name))
		return nil
	}

	if _, err := global.GetMod(name); err != nil {
		stats.AddWarning(fmt.Sprintf("cmd: Could find '%s'", name))
		return nil
	}

	if err := install.Install(name, true); err != nil {
		return fmt.Errorf("cmd: Could not update '%s' because an error occured in %s", name, err.Error())
	}

	stats.Updated++

	mod, err := local.GetMod(name)
	if err != nil {
		stats.AddWarning(fmt.Sprintf("cmd: Updated '%s' but it seems to be an invalid mod", name))
		return nil
	}

	return installDependencies(mod, stats)
}
