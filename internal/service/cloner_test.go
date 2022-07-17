package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	repo             = "https://github.com/octocat/Hello-World"
	branch           = "master"
	latestCommitHash = "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"
)

func TestCloner_GetLatestCommitHash(t *testing.T) {
	var cloner Cloner

	hash, err := cloner.GetLatestCommitHash(repo, branch)
	assert.NoError(t, err)
	assert.Equal(t, latestCommitHash, hash)

	hash, err = cloner.GetLatestCommitHash(repo, "-")
	assert.ErrorIs(t, err, ErrBranchNotFound)
	assert.Empty(t, hash)
}

func TestCloner_CloneRepository(t *testing.T) {
	var cloner = NewCloner(t.TempDir(), &TarArchiver{})

	archivePath, removeArchive, err := cloner.CloneRepository(repo, branch, latestCommitHash)
	assert.NoError(t, err)
	assert.FileExists(t, archivePath)

	removeArchive()
	assert.NoFileExists(t, archivePath)
}
