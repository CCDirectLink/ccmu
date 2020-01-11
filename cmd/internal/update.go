package internal

import (
	"errors"
	"fmt"

	"github.com/CCDirectLink/ccmu/pkg"
)

//Update a mod
func Update(args []string) {
	var all []pkg.Package
	if len(args) == 0 {
		all, _ = getGame().Installed()
	} else {
		all = getAll(args)
	}

	for _, p := range all {
		err := p.Update()
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
