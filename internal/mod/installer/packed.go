package installer

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

//Packed represents the installer for a packed mod which can either be local or remote.
type Packed struct {
	Path   string
	ModDir string
}

//Install the packed mod.
func (p Packed) Install() error {
	var (
		reader io.Reader
		name   string
		err    error
		file   *os.File
	)

	if ok, urlpath := isValidURL(p.Path); ok {
		name = urlpath
		reader, err = download(p.Path)
	} else {
		name = p.Path
		file, err = os.Open(p.Path)
		defer file.Close()
		reader = file
	}
	if err != nil {
		return err
	}

	name = filepath.Base(name)
	name = name[:len(name)-len(filepath.Ext(name))] + ".ccmod"

	target, err := os.Create(filepath.Join(p.ModDir, name))
	_, err = io.Copy(target, reader)
	return err
}

func isValidURL(toTest string) (bool, string) {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false, ""
	}

	u, err := url.Parse(toTest)
	return err == nil && u.Scheme != "" && u.Host != "", u.Path
}

func download(url string) (*bytes.Buffer, error) {
	resp, err := http.Get(url)
	if err != nil {
		return &bytes.Buffer{}, err
	}
	defer resp.Body.Close()

	var result bytes.Buffer
	_, err = io.Copy(&result, resp.Body)
	return &result, err
}
