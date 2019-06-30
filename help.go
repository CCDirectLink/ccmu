package main

import "fmt"

func printHelp() {
	fmt.Println("Usage: ccmu [command] [options] <args...>")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("")
	fmt.Println("install <mod name>")
	fmt.Println("uninstall <mod name>")
	fmt.Println("update <mod name>")
	fmt.Println("list")
	fmt.Println("outdated")
	fmt.Println("help")
}
