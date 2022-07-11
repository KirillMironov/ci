package service

import (
	"github.com/go-git/go-git/v5"
	"os"
)

// Cloner is a service that can clone a repository.
type Cloner struct{}

// CloneRepository clones a repository to a temporary directory and returns its path and a function that removes it.
func (Cloner) CloneRepository(url string) (sourceCodePath string, removeSourceCode func() error, err error) {
	sourceCodePath, err = os.MkdirTemp("", "")
	if err != nil {
		return "", nil, err
	}

	_, err = git.PlainClone(sourceCodePath, false, &git.CloneOptions{URL: url})
	if err != nil {
		return "", nil, err
	}

	return sourceCodePath, func() error { return os.RemoveAll(sourceCodePath) }, nil
}
