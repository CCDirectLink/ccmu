package pkg

//Package represents a common set of functions all mods and tools share.
type Package interface {
	Info() (Info, error)
	Installed() bool

	Install() error
	Uninstall() error
	Update() error

	Dependencies() ([]Package, error)
}
