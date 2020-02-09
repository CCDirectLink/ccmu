package mod

import (
	"path/filepath"

	"github.com/CCDirectLink/ccmu/internal/game"
)

//Mod package. Implements pkg.Package.
type Mod struct {
	Name   string
	Path   string
	Game   game.Game
	Source string

	realPath string
}

//FromSource creates a mod that installed from the specified source.
func FromSource(source string, game game.Game) (Mod, error) {
	basePath, _ := game.BasePath()
	pkg, err := readPackageFromSource(source)

	return Mod{
		Name:   pkg.Name,
		Source: source,
		Game:   game,
		Path:   basePath,
	}, err
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
