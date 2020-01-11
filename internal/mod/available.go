package mod

import "github.com/CCDirectLink/ccmu/internal/moddb"

//Available checks if it can be installed.
func (m Mod) Available() bool {
	_, err := moddb.PkgInfo(m.Name)
	return err == nil
}
