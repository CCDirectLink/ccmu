package mod

import (
	"os"

	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

//Uninstall the mod.
func (m Mod) Uninstall() error {
	path, err := m.local()
	if err != nil {
		return pkg.NewErrorReason(pkg.ReasonNotFound, pkg.ModeUninstall, m, err)
	}

	err = os.RemoveAll(path)
	if err != nil {
		return pkg.NewError(pkg.ModeUninstall, m, err)
	}
	return nil
}
