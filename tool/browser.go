package tool

import (
	"github.com/CCDirectLink/ccmu/internal/game"
	"github.com/CCDirectLink/ccmu/pkg"
	"github.com/CCDirectLink/ccmu/tool/registry"
	"os"
)

type browser struct {
	game game.Game
}

func (b browser) Info() (pkg.Info, error) {
	info := pkg.Info{
		Name:           "browser",
		NiceName:       "Browser Extension",
		Description:    "Install this extension to register an extension to call ccmu from the browser.",
		Licence:        "",
		CurrentVersion: "0.0.0",
		NewestVersion:  "0.0.0",
		Hidden:         false,
	}

	if b.Available() {
		info.NewestVersion = "1.0.0"
	}

	if b.Installed() {
		path, _ := b.game.BasePath()
		if path == registry.ProtocolInstalled() {
			info.CurrentVersion = "1.0.0"
		} else {
			info.CurrentVersion = "0.1.0"
		}
	}

	return info, nil
}

func (b browser) Installed() bool {
	return registry.ProtocolInstalled() != ""
}

func (b browser) Available() bool {
	return registry.Supported
}

func (b browser) Install() error {
	path, err := b.game.BasePath()
	if err != nil {
		return err
	}
	return registry.RegisterProtocol(os.Args[0], path)
}

func (b browser) Uninstall() error {
	return registry.UnregisterProtocol()
}

func (b browser) Update() error {
	err := b.Uninstall()
	if err != nil {
		return err
	}
	return b.Install()
}

func (b browser) Dependencies() ([]pkg.Package, error) {
	return []pkg.Package{}, nil
}

func (b browser) NewestDependencies() ([]pkg.Package, error) {
	return []pkg.Package{}, nil
}
