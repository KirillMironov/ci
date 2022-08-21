package service

import (
	"errors"
	"fmt"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"path/filepath"
)

var (
	ErrBranchNotFound     = errors.New("branch not found")
	ErrRepositoryNotFound = errors.New("repository not found")
	ErrRevisionNotFound   = errors.New("revision not found")
)

// Cloner used to clone source code repositories.
type Cloner struct {
	// Path to the directory where repositories are stored.
	repositoriesDir string
}

func NewCloner(repositoriesDir string) *Cloner {
	return &Cloner{repositoriesDir: repositoriesDir}
}

// GetLatestCommitHash returns the hash of the latest commit in the given repository branch.
func (Cloner) GetLatestCommitHash(repo domain.Repository) (string, error) {
	var remote = git.NewRemote(nil, &config.RemoteConfig{URLs: []string{repo.URL}})
	var targetReference = plumbing.NewBranchReferenceName(repo.Branch).String()

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		if errors.Is(err, transport.ErrRepositoryNotFound) {
			return "", ErrRepositoryNotFound
		}
		return "", err
	}

	for _, ref := range refs {
		if ref.Name().String() == targetReference {
			return ref.Hash().String(), nil
		}
	}

	return "", ErrBranchNotFound
}

// CloneRepository clones a repository and returns the local repository path.
func (c Cloner) CloneRepository(repo domain.Repository, targetHash string) (srcCodePath string, err error) {
	repository, srcCodePath, err := c.openOrCloneRepository(repo)
	if err != nil {
		return "", fmt.Errorf("failed to open or clone repository: %w", err)
	}

	wt, err := repository.Worktree()
	if err != nil {
		return "", err
	}

	err = wt.Pull(&git.PullOptions{
		ReferenceName: plumbing.NewBranchReferenceName(repo.Branch),
		SingleBranch:  true,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return "", err
	}

	revision, err := repository.ResolveRevision(plumbing.Revision(targetHash))
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return "", ErrRevisionNotFound
		}
		return "", err
	}

	return srcCodePath, wt.Checkout(&git.CheckoutOptions{Hash: *revision})
}

func (c Cloner) openOrCloneRepository(repo domain.Repository) (repository *git.Repository, localPath string, _ error) {
	abs, err := filepath.Abs(c.repositoriesDir)
	if err != nil {
		return nil, "", err
	}
	localPath = filepath.Join(abs, repo.Id)

	repository, err = git.PlainOpen(localPath)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			repository, err = git.PlainClone(localPath, false, &git.CloneOptions{
				URL:           repo.URL,
				ReferenceName: plumbing.NewBranchReferenceName(repo.Branch),
				NoCheckout:    true,
			})
			switch {
			case errors.Is(err, transport.ErrRepositoryNotFound):
				return nil, "", ErrRepositoryNotFound
			case errors.Is(err, plumbing.ErrReferenceNotFound):
				return nil, "", ErrBranchNotFound
			case err != nil:
				return nil, "", err
			}
		} else {
			return nil, "", err
		}
	}

	return repository, localPath, nil
}
