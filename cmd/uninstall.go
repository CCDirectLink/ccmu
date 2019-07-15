package cmd

import (
	"fmt"
	"os"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/local"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal/tools"
)

//Uninstall removes a mod from a directory
func Uninstall(args []string) (*Stats, error) {
	if _, err := local.GetGame(); err != nil {
		return nil, fmt.Errorf("cmd: Could not find game folder")
	}

	stats := &Stats{}
	for _, name := range args {
		mod, err := local.GetMod(name)
		if err != nil {
			err = uninstallTool(name, stats)
			if err != nil {
				return stats, err
			}
			continue
		}

		err = os.RemoveAll(mod.BasePath)
		if err != nil {
			stats.AddWarning(fmt.Sprintf("cmd: Could not remove mod '%s' because of an error in %s", name, err.Error()))
		}

		stats.Removed++
	}

	return stats, nil
}

func uninstallTool(name string, stats *Stats) error {
	tool := tools.Find(name)
	if tool == nil {
		stats.AddWarning(fmt.Sprintf("cmd: Could not find mod or tool '%s'", name))
		return nil
	}

	err := tool.Uninstall()
	if err != nil {
		return err
	}

	stats.Removed++
	return nil
}
