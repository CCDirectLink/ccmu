package tools

type ccloader struct{}

func (ccloader) Newest() (string, error) {
	return "0.0.0", nil
}
func (ccloader) Current() (string, error) {
	return "0.0.0", nil
}

func (ccloader) Install() error {
	return nil
}
func (ccloader) Uninstall() error {
	return nil
}
func (ccloader) Update() error {
	return nil
}
