package mod

import (
	"errors"
	"os"

	"github.com/CCDirectLink/ccmu/pkg"
)

//Uninstall the mod.
func (m Mod) Uninstall() error {
	path, err := m.local()
	if err != nil {
		return pkg.NewErrorReason(pkg.ReasonNotFound, pkg.ModeUninstall, m, err)
	}

	allMods, err := m.Game.Installed()
	if err != nil {
		var pkgErr pkg.Error
		if errors.As(err, &pkgErr) {
			pkgErr.Mode = pkg.ModeUninstall
			return pkgErr
		}
		return err
	}

	for _, mod := range allMods {
		deps, _ := mod.Dependencies()
		for _, rep := range deps {
			info, _ := rep.Info()
			if info.Name == m.Name {
				return pkg.NewErrorReason(pkg.ReasonDependant, pkg.ModeUninstall, m, nil)
			}
		}
	}

	err = os.RemoveAll(path)
	if err != nil {
		return pkg.NewError(pkg.ModeUninstall, m, err)
	}
	return nil
}
