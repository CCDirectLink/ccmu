package internal

import (
	"errors"
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

var installed = 0

//Install a mod
func Install(args []string) {
	for _, p := range getAll(args) {
		err := p.Install()
		if err != nil {
			var pkgErr pkg.Error
			if errors.As(err, &pkgErr) {
				fmt.Println(pkgErr.String())
			} else {
				fmt.Printf("An error occured in %s\n", err)
			}
			continue
		}
	}
}
