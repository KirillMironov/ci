package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloner_GetLatestCommitHash(t *testing.T) {
	var cloner Cloner

	hash, err := cloner.GetLatestCommitHash("https://github.com/KirillMironov/ci", "main")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	hash, err = cloner.GetLatestCommitHash("https://github.com/KirillMironov/ci", "-")
	assert.ErrorIs(t, err, ErrBranchNotFound)
	assert.Empty(t, hash)
}

func TestCloner_CloneRepository(t *testing.T) {
	var cloner Cloner

	sourceCodePath, removeSourceCode, err := cloner.CloneRepository("https://github.com/KirillMironov/ci",
		"main", "afa50416019b2583da7d7f1e6ae26a511273031e")
	assert.NoError(t, err)
	assert.NotEmpty(t, sourceCodePath)
	assert.DirExists(t, sourceCodePath)

	removeSourceCode()
	assert.NoDirExists(t, sourceCodePath)
}
