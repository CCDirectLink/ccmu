package pkg

import "github.com/Masterminds/semver"

//Info about a package.
type Info struct {
	Name           string
	NiceName       string
	Description    string
	Licence        string
	CurrentVersion string
	NewestVersion  string
	Hidden         bool
}

//Outdated checks if an newer version is available
func (info Info) Outdated() (bool, error) {
	newest, err := semver.NewVersion(info.NewestVersion)
	if err != nil {
		return false, NewError(ModeUnknown, nil, err)
	}

	current, err := semver.NewVersion(info.CurrentVersion)
	if err != nil {
		return false, NewError(ModeUnknown, nil, err)
	}

	return current.LessThan(newest), nil
}
