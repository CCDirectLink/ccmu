package mod

//Installed checks if the mod exists or not by checking if it's package.json exists.
func (m Mod) Installed() bool {
	_, err := getMod(m.Name, m.Path, m.Game)
	return err == nil
}
