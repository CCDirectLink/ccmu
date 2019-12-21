package mod

import "github.com/CCDirectLink/CCUpdaterCLI/pkg"

//Dependencies of a mod including indirect ones.
func (m Mod) Dependencies() ([]pkg.Package, error) {
	all, err := m.allDeps()

	var result []pkg.Package
	for _, pkg := range all {
		pkgInfo, _ := pkg.Info()

		duplicate := false
		for _, p := range result {
			info, _ := p.Info()
			if pkgInfo.Name == info.Name {
				duplicate = true
				break
			}
		}

		if !duplicate {
			result = append(result, pkg)
		}
	}

	return result, err
}

func (m Mod) directDeps() ([]pkg.Package, error) {
	data, err := m.readPackageFile()
	if err != nil {
		return []pkg.Package{}, err
	}

	var result = make([]pkg.Package, len(data.Dependencies))
	var i = 0
	for name := range data.Dependencies {
		result[i], err = m.Game.Get(name)
		if err != nil {
			return result, err
		}
		i++
	}

	return result, nil
}

func (m Mod) allDeps() ([]pkg.Package, error) {
	var result []pkg.Package

	direct, err := m.directDeps()
	if err != nil {
		return direct, err
	}

	for _, pkg := range direct {
		indirect, err := pkg.Dependencies()
		if err != nil {
			return result, err
		}

		result = append(append(result, pkg), indirect...)
	}

	return result, nil
}
