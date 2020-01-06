package mod

import "github.com/CCDirectLink/CCUpdaterCLI/pkg"

//Dependencies of a mod including indirect ones.
func (m Mod) Dependencies() ([]pkg.Package, error) {
	all, err := m.allLocalDeps()
	return removeDuplicates(all), err
}

func (m Mod) directLocalDeps() ([]pkg.Package, error) {
	data, err := m.readPackageFile()
	if err != nil {
		return []pkg.Package{}, err
	}

	return m.mapDeps(data.Dependencies)
}

func (m Mod) allLocalDeps() ([]pkg.Package, error) {
	var result []pkg.Package

	direct, err := m.directDeps()
	if err != nil {
		return direct, err
	}

	result = append(result, direct...)

	for _, pkg := range direct {
		indirect, err := pkg.Dependencies()
		if err != nil {
			return result, err
		}

		result = append(result, indirect...)
	}

	return result, nil
}
