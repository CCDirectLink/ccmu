package mod

import (
	"errors"
	"path/filepath"

	"github.com/CCDirectLink/ccmu/internal/mod/installer"
	"github.com/CCDirectLink/ccmu/pkg"
)

//Install the mod
func (m Mod) Install() error {
	if m.Installed() {
		return pkg.NewErrorReason(pkg.ReasonAlreadyInstalled, pkg.ModeInstall, m, nil)
	}

	if !m.Available() {
		return pkg.NewErrorReason(pkg.ReasonNotFound, pkg.ModeInstall, m, nil)
	}

	err := m.install()
	if err != nil {
		return pkg.NewError(pkg.ModeInstall, m, err)
	}
	return err
}

func (m Mod) install() error {
	deps, err := m.directDeps()
	if err != nil {
		return err
	}

	for _, dep := range deps {
		err = m.installDep(dep)
		if err != nil {
			return err
		}
	}

	return m.installTo(m.path())

}

func (m Mod) installTo(path string) error {
	var inst installer.Installer
	if m.Source == "" {
		inst = installer.Moddb{
			Name:         m.Name,
			Path:         path,
			PreferPacked: true,
		}
	} else {
		inst = installer.Packed{
			Path:   m.Source,
			ModDir: filepath.Dir(path),
		}
	}

	return inst.Install()
}

func (m Mod) installDep(dep pkg.Package) error {
	var err error
	if !dep.Installed() {
		err = dep.Install()
	} else {
		info, err := dep.Info()
		if err == nil {
			if outdated, _ := info.Outdated(); outdated {
				err = dep.Update()
			}
		}
	}

	//Do not return error when already installed
	var pkgErr pkg.Error
	if errors.As(err, &pkgErr) && pkgErr.Reason == pkg.ReasonAlreadyInstalled {
		return nil
	}
	return err
}
