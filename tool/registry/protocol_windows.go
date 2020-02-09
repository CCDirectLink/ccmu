package registry

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

//Supported shows if registry is supported or not.
const Supported = true

//RegisterProtocol in HKEY_CLASSES_ROOT.
func RegisterProtocol(ccmu, game string) error {
	software, _ := registry.OpenKey(registry.CURRENT_USER, "Software", registry.READ|registry.WRITE)
	defer software.Close()
	classes, _ := registry.OpenKey(software, "Classes", registry.READ|registry.WRITE)
	defer software.Close()

	root, _, err := registry.CreateKey(classes, "ccmu", registry.READ|registry.WRITE)
	if err != nil {
		return err
	}
	defer root.Close()
	root.SetStringValue("", "URL:CrossCode Mod Updater Protocol")
	root.SetStringValue("URL Protocol", "")

	shell, _, _ := registry.CreateKey(root, "shell", registry.READ|registry.WRITE)
	defer shell.Close()
	open, _, _ := registry.CreateKey(shell, "open", registry.READ|registry.WRITE)
	defer open.Close()
	cmd, _, _ := registry.CreateKey(open, "command", registry.READ|registry.WRITE)
	defer cmd.Close()

	return cmd.SetStringValue("", fmt.Sprintf("\"%s\" -game=\"%s\" -url=\"%%1\" install", ccmu, strings.ReplaceAll(game, "\\", "\\\\")))
}

//UnregisterProtocol from HKEY_CLASSES_ROOT.
func UnregisterProtocol() error {
	software, _ := registry.OpenKey(registry.CURRENT_USER, "Software", registry.WRITE)
	defer software.Close()
	classes, _ := registry.OpenKey(software, "Classes", registry.WRITE)
	defer software.Close()
	return registry.DeleteKey(classes, "ccmu")
}

//ProtocolInstalled in HKEY_CLASSES_ROOT.
func ProtocolInstalled() string {
	software, _ := registry.OpenKey(registry.CURRENT_USER, "Software", registry.WRITE)
	defer software.Close()
	classes, _ := registry.OpenKey(software, "Classes", registry.WRITE)
	defer software.Close()

	root, err := registry.OpenKey(classes, "ccmu", registry.READ)
	defer root.Close()
	if err != nil {
		return ""
	}

	shell, _ := registry.OpenKey(root, "shell", registry.READ)
	defer shell.Close()
	open, _ := registry.OpenKey(shell, "open", registry.READ)
	defer open.Close()
	cmd, _ := registry.OpenKey(open, "command", registry.READ)
	defer cmd.Close()
	result, _, _ := cmd.GetStringValue("")
	return result
}
