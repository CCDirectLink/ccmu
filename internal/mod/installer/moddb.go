package installer

import (
	"archive/zip"
	"bytes"
	"os"

	"github.com/CCDirectLink/ccmu/internal/moddb"
)

//Moddb installs a mod with the given name.
type Moddb struct {
	Name         string
	Path         string
	PreferPacked bool
}

//Install a moddb mod.
func (m Moddb) Install() error {
	var (
		buf bytes.Buffer
		src string
	)

	buf, src, err := moddb.DownloadMod(m.Name, m.PreferPacked)
	if err != nil {
		return err
	}

	if src == "" {
		file, _ := os.Create(m.Path + ".ccmod")
		defer file.Close()
		_, err = file.Write(buf.Bytes())
	} else {
		zipReader, _ := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		_, err = extract(zipReader, m.Path, src)
	}

	return err
}
