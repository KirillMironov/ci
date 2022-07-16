package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArchiver(t *testing.T) {
	var archiver Archiver

	archivePath, removeArchive, err := archiver.Compress(".")
	assert.NoError(t, err)
	assert.NotEmpty(t, archivePath)
	assert.FileExists(t, archivePath)

	data, err := archiver.FindFile("archiver.go", archivePath)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	removeArchive()
	assert.NoFileExists(t, archivePath)
}
