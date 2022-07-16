package service

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

var ErrBranchNotFound = errors.New("branch not found")

// Cloner is a service that can clone a repository.
type Cloner struct{}

// GetLatestCommitHash returns the hash of the latest commit in the given repository branch.
func (Cloner) GetLatestCommitHash(url, branch string) (string, error) {
	remote := git.NewRemote(nil, &config.RemoteConfig{Name: "origin", URLs: []string{url}})

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return "", err
	}

	target := "refs/heads/" + branch

	for _, ref := range refs {
		if ref.Name().String() == target {
			return ref.Hash().String(), nil
		}
	}

	return "", ErrBranchNotFound
}

// CloneRepository clones a repository to a temporary directory and returns its path and a function that removes it.
func (Cloner) CloneRepository(url, branch, hash string) (path string, remove func(), err error) {
	path, err = os.MkdirTemp("", "")
	if err != nil {
		return "", nil, err
	}

	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           url,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		return "", nil, err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return "", nil, err
	}

	err = wt.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(hash)})
	if err != nil {
		return "", nil, err
	}

	return path, func() { os.RemoveAll(path) }, nil
}
