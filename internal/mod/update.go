package mod

//Update a mod if it's outdated.
func (m Mod) Update() error {
	info, err := m.Info()
	if err != nil {
		return err
	}

	if outdated, err := info.Outdated(); err != nil || !outdated {
		return err
	}

	err = m.Uninstall()
	if err != nil {
		return err
	}
	return m.Install()
}
