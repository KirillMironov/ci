package service

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	url              = "https://github.com/octocat/Hello-World"
	branch           = "master"
	latestCommitHash = "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"
)

func TestCloner_GetLatestCommitHash(t *testing.T) {
	var cloner Cloner

	tests := map[string]struct {
		repo               domain.Repository
		expectedCommitHash string
		expectedError      error
	}{
		"success": {
			repo:               domain.Repository{URL: url, Branch: branch},
			expectedCommitHash: latestCommitHash,
			expectedError:      nil,
		},
		"branch not found": {
			repo:               domain.Repository{URL: url, Branch: "-"},
			expectedCommitHash: "",
			expectedError:      ErrBranchNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			commitHash, err := cloner.GetLatestCommitHash(tc.repo)

			assert.Equal(t, tc.expectedCommitHash, commitHash)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCloner_CloneRepository(t *testing.T) {
	var cloner = NewCloner(t.TempDir(), &TarArchiver{})

	tests := map[string]struct {
		repo              domain.Repository
		targetCommitHash  string
		expectArchivePath bool
		expectError       bool
	}{
		"success": {
			repo:              domain.Repository{URL: url, Branch: branch},
			targetCommitHash:  latestCommitHash,
			expectArchivePath: true,
			expectError:       false,
		},
		"target commit not found": {
			repo:              domain.Repository{URL: url, Branch: "-"},
			targetCommitHash:  "-",
			expectArchivePath: false,
			expectError:       true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			archivePath, removeArchive, err := cloner.CloneRepository(tc.repo, tc.targetCommitHash)

			assert.Equal(t, tc.expectArchivePath, archivePath != "")
			assert.Equal(t, tc.expectError, err != nil)

			if tc.expectArchivePath {
				assert.FileExists(t, archivePath)
				removeArchive()
				assert.NoFileExists(t, archivePath)
			}
		})
	}
}
