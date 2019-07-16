package tools

type simplify struct {
	loader ccloader
}

func (simplify) Newest() (string, error) {
	return "0.0.0", nil
}
func (simplify) Current() (string, error) {
	return "0.0.0", nil
}

func (s simplify) Install() error {
	return s.loader.Install()
}
func (s simplify) Uninstall() error {
	return s.loader.Uninstall()
}
func (s simplify) Update() error {
	return s.loader.Update()
}
