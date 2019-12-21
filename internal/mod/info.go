package mod

import (
	"github.com/CCDirectLink/CCUpdaterCLI/internal/moddb"
	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

//Info about the mod.
func (m Mod) Info() (pkg.Info, error) {
	if !m.Installed() {
		return moddb.ModInfo(m.Name)
	}

	info := m.localInfo()
	err := moddb.MergeModInfo(&info)
	return info, err
}

func (m Mod) localInfo() pkg.Info {
	data, _ := m.readPackageFile()
	return pkg.Info{
		Name:           m.Name,
		NiceName:       data.Name,
		Description:    data.Description,
		Licence:        data.Licence,
		CurrentVersion: data.Version,
	}
}
