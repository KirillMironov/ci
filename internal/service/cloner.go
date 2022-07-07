package service

import (
	"os"
	"os/exec"
)

type Cloner struct {
	repositoriesPath string
}

func NewCloner(repositoriesPath string) *Cloner {
	return &Cloner{repositoriesPath: repositoriesPath}
}

func (c Cloner) CloneRepository(url string) (dir string, err error) {
	dir, err = os.MkdirTemp(c.repositoriesPath, "")
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
