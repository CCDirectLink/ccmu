package internal

import (
	"errors"
	"fmt"

	"github.com/CCDirectLink/ccmu/pkg"
)

//Uninstall removes a mod from a directory
func Uninstall(args []string) {
	for _, p := range getAll(args) {
		err := p.Uninstall()
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
