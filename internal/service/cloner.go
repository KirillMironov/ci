package service

import (
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
	"path/filepath"
)

var ErrBranchNotFound = errors.New("branch not found")

// Cloner used to clone source code repositories.
type Cloner struct {
	repositoriesDir string // Path to the directory where repositories are stored.
	archiver        archiver
}

type archiver interface {
	Compress(dir string) (archivePath string, err error)
}

func NewCloner(repositoriesDir string, archiver archiver) *Cloner {
	return &Cloner{
		repositoriesDir: repositoriesDir,
		archiver:        archiver,
	}
}

// GetLatestCommitHash returns the hash of the latest commit in the given repository branch.
func (Cloner) GetLatestCommitHash(repo domain.Repository) (string, error) {
	remote := git.NewRemote(nil, &config.RemoteConfig{Name: "origin", URLs: []string{repo.URL}})

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return "", err
	}

	target := "refs/heads/" + repo.Branch

	for _, ref := range refs {
		if ref.Name().String() == target {
			return ref.Hash().String(), nil
		}
	}

	return "", ErrBranchNotFound
}

// CloneRepository clones a repository and returns the path to the compressed source code archive.
func (c Cloner) CloneRepository(repo domain.Repository, hash string) (archivePath string, removeArchive func(),
	err error) {
	abs, err := filepath.Abs(c.repositoriesDir)
	if err != nil {
		return "", nil, err
	}

	localPath := filepath.Join(abs, repo.Id)

	var repository *git.Repository

	repository, err = git.PlainOpen(localPath)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			repository, err = git.PlainClone(localPath, false, &git.CloneOptions{
				URL:           repo.URL,
				SingleBranch:  true,
				ReferenceName: plumbing.NewBranchReferenceName(repo.Branch),
			})
			if err != nil {
				return "", nil, err
			}
		} else {
			return "", nil, err
		}
	}

	wt, err := repository.Worktree()
	if err != nil {
		return "", nil, err
	}

	err = wt.Pull(&git.PullOptions{
		ReferenceName: plumbing.NewBranchReferenceName(repo.Branch),
		SingleBranch:  true,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return "", nil, err
	}

	err = wt.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(hash)})
	if err != nil {
		return "", nil, err
	}

	archivePath, err = c.archiver.Compress(localPath)
	if err != nil {
		return "", nil, err
	}

	return archivePath, func() { os.Remove(archivePath) }, nil
}
