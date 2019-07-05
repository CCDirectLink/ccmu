package main

import "fmt"

func printHelp() {
	fmt.Println("Usage: ccmu [command] [options] <args...>")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --game <path>         Sets the game folder used for operations")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  install <mod name>    Installs one or more mods")
	fmt.Println("  uninstall <mod name>  Uninstall one or more mods")
	fmt.Println("  update <mod name>     Updates one or more mods")
	fmt.Println("  list                  Lists all mods that the tool knows about")
	fmt.Println("  outdated              Show the names and versions of outdated mods")
	fmt.Println("  version               Display the version of this tool")
	fmt.Println("  help                  Display this message")
}
