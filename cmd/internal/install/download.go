package install

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func download(url string) (*os.File, error) {
	file, err := ioutil.TempFile("installing", "mod")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return file, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return file, err
}
