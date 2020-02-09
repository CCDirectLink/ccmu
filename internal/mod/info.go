package mod

import (
	"github.com/CCDirectLink/ccmu/internal/moddb"
	"github.com/CCDirectLink/ccmu/pkg"
)

//Info about the mod.
func (m Mod) Info() (pkg.Info, error) {

	if !m.Installed() && m.Source == "" {
		if !m.Available() {
			return pkg.Info{
				Name:     m.Name,
				NiceName: m.Name,
			}, pkg.ErrNotFound
		}

		return moddb.PkgInfo(m.Name)
	}

	info := m.localInfo()
	var err error
	if m.Available() && m.Source == "" {
		err = moddb.MergePkgInfo(&info)
	}

	if m.Source != "" {
		data, _ := readPackageFromSource(m.Source)
		info.Name = data.Name
		info.NiceName = data.Name
		info.Description = data.Description
		info.Licence = data.Licence
		info.NewestVersion = data.Version
		info.Hidden = data.Hidden
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
		Hidden:         data.Hidden,
	}
}
