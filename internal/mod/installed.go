package mod

import (
	"os"
	"path/filepath"
)

//Installed checks if the mod exists or not by checking if it's package.json exists.
func (m Mod) Installed() bool {
	path := filepath.Join(m.path(), "package.json")
	_, err := os.Stat(path)
	return err == nil
}
