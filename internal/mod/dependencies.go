package mod

import (
	"github.com/CCDirectLink/CCUpdaterCLI/internal/moddb"
	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

//NewestDependencies of a mod including indirect ones.
func (m Mod) NewestDependencies() ([]pkg.Package, error) {
	all, err := m.allDeps()
	return removeDuplicates(all), err
}

func removeDuplicates(all []pkg.Package) []pkg.Package {
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
	return result
}

func (m Mod) mapDeps(data map[string]string) ([]pkg.Package, error) {
	var err error
	result := make([]pkg.Package, len(data))
	i := 0
	for name := range data {
		result[i], err = m.Game.Get(name)
		if err != nil {
			return result, err
		}
		i++
	}

	return result, nil
}

func (m Mod) directDeps() ([]pkg.Package, error) {
	data, err := moddb.Dependencies(m.Name)
	if err != nil {
		return []pkg.Package{}, err
	}
	return m.mapDeps(data)
}

func (m Mod) allDeps() ([]pkg.Package, error) {
	var result []pkg.Package

	direct, err := m.directDeps()
	if err != nil {
		return direct, err
	}

	result = append(result, direct...)

	for _, pkg := range direct {
		indirect, err := pkg.NewestDependencies()
		if err != nil {
			return result, err
		}

		result = append(result, indirect...)
	}

	return result, nil
}
