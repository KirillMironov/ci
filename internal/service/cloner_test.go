package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloner_CloneRepository(t *testing.T) {
	var cloner Cloner

	sourceCodePath, removeSourceCode, err := cloner.CloneRepository("https://github.com/KirillMironov/kube")
	assert.NoError(t, err)
	assert.NotEmpty(t, sourceCodePath)
	assert.DirExists(t, sourceCodePath)

	assert.NoError(t, removeSourceCode())
	assert.NoDirExists(t, sourceCodePath)
}
