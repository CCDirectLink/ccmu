package mod

import "github.com/CCDirectLink/CCUpdaterCLI/pkg"

//Update a mod if it's outdated.
func (m Mod) Update() error {
	info, err := m.Info()
	if err != nil {
		if pkgErr, ok := err.(pkg.Error); ok {
			pkgErr.Mode = pkg.ModeUpdate
			return pkgErr
		}
		return err
	}

	if outdated, err := info.Outdated(); err != nil || !outdated {
		if pkgErr, ok := err.(pkg.Error); ok {
			pkgErr.Mode = pkg.ModeUpdate
			return pkgErr
		}
		return err
	}

	err = m.Uninstall()
	if err != nil {
		if pkgErr, ok := err.(pkg.Error); ok {
			pkgErr.Mode = pkg.ModeUpdate
			return pkgErr
		}
		return err
	}
	err = m.Install()
	if err != nil {
		if pkgErr, ok := err.(pkg.Error); ok {
			pkgErr.Mode = pkg.ModeUpdate
			return pkgErr
		}
		return err
	}
	return nil
}
