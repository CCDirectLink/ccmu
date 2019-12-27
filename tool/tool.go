package tool

import (
	"strings"

	"github.com/CCDirectLink/CCUpdaterCLI/internal/moddb"
	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

//Available tools.
func Available(path string) ([]pkg.Package, error) {
	return []pkg.Package{
		ccloader{path},
	}, nil
}

//Installed tools.
func Installed(path string) ([]pkg.Package, error) {
	avail, err := Available(path)
	result := make([]pkg.Package, 0, len(avail))

	for _, tool := range avail {
		if tool.Installed() {
			result = append(result, tool)
		}
	}

	return result, err
}

//Get a tool by exact name.
func Get(path, name string) (pkg.Package, error) {
	switch strings.ToLower(name) {
	case "ccloader":
		return ccloader{path}, nil
	default:
		return nil, pkg.NewError(pkg.ModeUnknown, nil, moddb.ErrNotFound)
	}
}
