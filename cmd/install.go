package cmd

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/global"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/install"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/tools"
)

var installed = 0

//Install a mod
func Install(args []string) (*Stats, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("cmd: No mods installed since no mods were specified")
	}

	if _, err := local.GetGame(); err != nil {
		return nil, fmt.Errorf("cmd: Could not find game folder")
	}

	if _, err := global.FetchModData(); err != nil {
		return nil, fmt.Errorf("cmd: Could not download mod data because an error occured in %s", err.Error())
	}

	stats := &Stats{}

	for _, name := range args {
		if _, err := local.GetMod(name); err == nil {
			stats.AddWarning(fmt.Sprintf("cmd: Could not install '%s' because it was already installed", name))
			continue
		}

		if err := installMod(name, stats); err != nil {
			return stats, err
		}
	}

	return stats, nil
}

func installMod(name string, stats *Stats) error {
	if _, err := global.GetMod(name); err != nil {
		return installTool(name, stats)
	}

	if err := install.Install(name, false); err != nil {
		return fmt.Errorf("cmd: Could not install '%s' because an error occured in %s", name, err.Error())
	}

	mod, err := local.GetMod(name)
	if err != nil {
		stats.AddWarning(fmt.Sprintf("cmd: Installed '%s' but it seems to be an invalid mod", name))
		return nil
	}

	stats.Installed++
	return installDependencies(mod, stats)
}

func installTool(name string, stats *Stats) error {
	tool := tools.Find(name)
	if tool == nil {
		stats.AddWarning(fmt.Sprintf("cmd: Could find mod or tool '%s'", name))
		return nil
	}

	err := tool.Install()
	if err != nil {
		return err
	}

	stats.Installed++
	return nil
}
