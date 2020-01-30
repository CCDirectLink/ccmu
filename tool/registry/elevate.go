package registry

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func elevate() error {
	if !isAdmin() {
		fmt.Printf("Administrator permission required. Restart as administrator or enter passord:\n")

		host, err := os.Hostname()
		if err != nil {
			return err
		}

		args := "\"" + strings.Join(os.Args, "\" \"") + "\""
		cmd := exec.Command("runas", fmt.Sprintf("/user:%s\\administrator", host), args)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			return err
		}
		os.Exit(0)
	}
	return nil
}

func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}
