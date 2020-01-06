package mod

import (
	"path/filepath"

	"github.com/CCDirectLink/CCUpdaterCLI/internal/game"
)

//Mod package. Implements pkg.Package.
type Mod struct {
	Name string
	Path string
	Game game.Game

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
