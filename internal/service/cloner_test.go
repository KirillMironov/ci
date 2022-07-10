package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloner_Clone(t *testing.T) {
	var cloner Cloner

	archive, remove, err := cloner.CloneRepository("https://github.com/KirillMironov/kube")
	assert.NoError(t, err)
	assert.NotNil(t, archive)
	assert.NoError(t, remove())
}
