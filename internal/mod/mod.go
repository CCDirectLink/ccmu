package mod

import (
	"path/filepath"

	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

//game is an interface to avoid cyclic imports with ccmu.Game.
type game interface {
	Get(name string) (pkg.Package, error)
}

//Mod package. Implements pkg.Package.
type Mod struct {
	Name string
	Path string
	Game game

	realPath string
}

func (m Mod) path() string {
	if m.realPath != "" {
		return m.realPath
	}

	if !m.Installed() {
		m.realPath = filepath.Join(m.Path, "assets", "mods", m.Name)
		return m.realPath
	}

	path, _ := m.local()
	return path
}
