package tools

type crosscode struct{}

func (crosscode) Newest() (string, error) {
	return "0.0.0", nil
}
func (crosscode) Current() (string, error) {
	return "0.0.0", nil
}

func (crosscode) Install() error {
	return nil
}
func (crosscode) Uninstall() error {
	return nil
}
func (crosscode) Update() error {
	return nil
}
