package mod

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

//Package represents the data inside package.json
type Package struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Version      string            `json:"version"`
	Licence      string            `json:"licence"`
	Dependencies map[string]string `json:"ccmodDependecies"`
	Hidden       bool              `json:"hidden"`
}

func readPackage(reader io.Reader) (Package, error) {
	var result Package
	err := json.NewDecoder(reader).Decode(&result)
	return result, err
}

func (m Mod) readPackageFile() (Package, error) {
	path := filepath.Join(m.path(), "package.json")
	file, err := os.Open(path)
	if err != nil {
		return Package{}, nil
	}
	defer file.Close()

	return readPackage(file)
}
