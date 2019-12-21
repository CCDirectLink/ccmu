package mod

import "os"

//Uninstall the mod.
func (m Mod) Uninstall() error {
	return os.RemoveAll(m.path())
}
