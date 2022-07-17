package service

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTarArchiver(t *testing.T) {
	var archiver TarArchiver

	archivePath, err := archiver.Compress(".")
	assert.NoError(t, err)
	assert.NotEmpty(t, archivePath)
	assert.FileExists(t, archivePath)
	defer os.Remove(archivePath)

	data, err := archiver.FindFile("tar_archiver.go", archivePath)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	data, err = archiver.FindFile("-", archivePath)
	assert.Error(t, err)
	assert.Empty(t, data)
}
