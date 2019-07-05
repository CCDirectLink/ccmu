package install

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

func copyDir(dst, src string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcStat.Mode())
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcFile := path.Join(src, file.Name())
		dstFile := path.Join(dst, file.Name())

		if file.IsDir() {
			err = copyDir(dstFile, srcFile)
		} else {
			err = copyFile(dstFile, srcFile)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(dst, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcStat.Mode())
}
