package installer

//Installer exposes a simple install function that can be used to install a mod.
type Installer interface {
	Install() error
}
