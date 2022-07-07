package service

import (
	"io/ioutil"
	"os"
	"os/exec"
)

type Cloner struct{}

func (Cloner) Clone(url string) (dir string, err error) {
	dir, err = ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		return "", err
	}

	err = exec.Command("git", "clone", url, dir).Run()
	if err != nil {
		os.RemoveAll(dir)
		return "", err
	}

	return dir, nil
}
