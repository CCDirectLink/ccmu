package tools

import (
	"github.com/Masterminds/semver"
)

//Tool defines an interface that allows implementation of multiple tools
type Tool interface {
	Newest() (string, error)
	Current() (string, error)

	Install() error
	Uninstall() error
	Update() error
}

//Find tool by the given name. Returns nil if no tool was found
func Find(name string) Tool {
	switch name {
	case "ccloader":
		return &ccloader{}
	case "Simplify":
		return &simplify{}
	case "crosscode":
		return &crosscode{}
	default:
		return nil
	}
}

//Outdated checks if an newer version is available
func Outdated(tool Tool) (bool, error) {
	new, err := tool.Newest()
	if err != nil {
		return false, err
	}
	cur, err := tool.Current()
	if err != nil {
		return false, err
	}

	newest, err := semver.NewVersion(new)
	if err != nil {
		return false, err
	}
	current, err := semver.NewVersion(cur)
	if err != nil {
		return false, err
	}

	return current.LessThan(newest), nil
}
