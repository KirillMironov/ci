package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloner_Clone(t *testing.T) {
	var cloner Cloner

	arch, err := cloner.CloneRepository("https://github.com/KirillMironov/kube")
	assert.NoError(t, err)
	assert.NotNil(t, arch)
	arch.Close()
}
