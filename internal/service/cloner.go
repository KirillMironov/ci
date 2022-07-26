package service

import (
	"encoding/hex"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"hash/fnv"
	"os"
	"path/filepath"
)

var ErrBranchNotFound = errors.New("branch not found")

// Cloner is a service that can clone a repository.
type Cloner struct {
	reposDir string
	archiver archiver
}

// archiver is a service that can archive a directory.
type archiver interface {
	Compress(dir string) (archivePath string, err error)
}

// NewCloner creates a new Cloner.
func NewCloner(reposDir string, archiver archiver) *Cloner {
	return &Cloner{
		reposDir: reposDir,
		archiver: archiver,
	}
}

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

// CloneRepository clones a repository and returns the path to the compressed source code.
func (c Cloner) CloneRepository(url, branch, hash string) (archivePath string, removeArchive func(), err error) {
	abs, err := filepath.Abs(c.reposDir)
	if err != nil {
		return "", nil, err
	}

	repoPath := filepath.Join(abs, hex.EncodeToString(fnv.New32().Sum([]byte(url))))

	var repo *git.Repository

	repo, err = git.PlainOpen(repoPath)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			repo, err = git.PlainClone(repoPath, false, &git.CloneOptions{
				URL:           url,
				SingleBranch:  true,
				ReferenceName: plumbing.NewBranchReferenceName(branch),
			})
			if err != nil {
				return "", nil, err
			}
		} else {
			return "", nil, err
		}
	}

	wt, err := repo.Worktree()
	if err != nil {
		return "", nil, err
	}

	err = wt.Pull(&git.PullOptions{
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return "", nil, err
	}

	err = wt.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(hash)})
	if err != nil {
		return "", nil, err
	}

	archivePath, err = c.archiver.Compress(repoPath)
	if err != nil {
		return "", nil, err
	}

	return archivePath, func() { os.Remove(archivePath) }, nil
}
