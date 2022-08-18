package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTarArchiver(t *testing.T) {
	var archiver TarArchiver

	archivePath, removeArchive, err := archiver.Compress(".")
	assert.NoError(t, err)

	assert.FileExists(t, archivePath)
	removeArchive()
	assert.NoFileExists(t, archivePath)
}
