package game

import "github.com/CCDirectLink/ccmu/pkg"

//Game is an interface to avoid cyclic imports with game.Game.
type Game interface {
	Get(name string) (pkg.Package, error)
	Installed() ([]pkg.Package, error)
	BasePath() (string, error)
}
