package mod

import (
	"github.com/CCDirectLink/CCUpdaterCLI/internal/moddb"
	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

//Info about the mod.
func (m Mod) Info() (pkg.Info, error) {
	if !m.Installed() {
		if !m.Available() {
			return pkg.Info{
				Name: m.Name,
				NiceName: m.Name,
			}, moddb.ErrNotFound
		}

		return moddb.PkgInfo(m.Name)
	}

	info := m.localInfo()

	var err error
	if m.Available() {
		err = moddb.MergePkgInfo(&info)
	}
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
