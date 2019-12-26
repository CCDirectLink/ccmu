package mod

//Installed checks if the mod exists or not by checking if it's package.json exists.
func (m Mod) Installed() bool {
	_, err := m.local()
	return err == nil
}
