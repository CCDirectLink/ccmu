package mod

import "github.com/CCDirectLink/CCUpdaterCLI/internal/moddb"

//Available checks if it can be installed.
func (m Mod) Available() bool {
	_, err := moddb.PkgInfo(m.Name)
	return err == nil
}
