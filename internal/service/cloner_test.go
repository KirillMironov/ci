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
	tests := map[string]struct {
		repo               domain.Repository
		expectedCommitHash string
		expectedError      error
	}{
		"success": {
			repo:               domain.Repository{Id: "0", URL: url, Branch: branch},
			expectedCommitHash: latestCommitHash,
			expectedError:      nil,
		},
		"branch not found": {
			repo:               domain.Repository{Id: "0", URL: url, Branch: "-"},
			expectedCommitHash: "",
			expectedError:      ErrBranchNotFound,
		},
		"repository not found": {
			repo:               domain.Repository{Id: "0", URL: "example.com", Branch: "main"},
			expectedCommitHash: "",
			expectedError:      ErrRepositoryNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var cloner = NewCloner(t.TempDir(), &TarArchiver{})

			commitHash, err := cloner.GetLatestCommitHash(tc.repo)

			assert.Equal(t, tc.expectedCommitHash, commitHash)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestCloner_CloneRepository(t *testing.T) {
	tests := map[string]struct {
		repo          domain.Repository
		targetHash    string
		expectedError error
	}{
		"success": {
			repo:          domain.Repository{Id: "0", URL: url, Branch: branch},
			targetHash:    latestCommitHash,
			expectedError: nil,
		},
		"revision not found": {
			repo:          domain.Repository{Id: "0", URL: url, Branch: branch},
			targetHash:    "-",
			expectedError: ErrRevisionNotFound,
		},
		"branch not found": {
			repo:          domain.Repository{Id: "0", URL: url, Branch: "-"},
			targetHash:    latestCommitHash,
			expectedError: ErrBranchNotFound,
		},
		"repository not found": {
			repo:          domain.Repository{Id: "0", URL: "example.com", Branch: "main"},
			targetHash:    latestCommitHash,
			expectedError: ErrRepositoryNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var cloner = NewCloner(t.TempDir(), &TarArchiver{})

			archivePath, removeArchive, err := cloner.CloneRepository(tc.repo, tc.targetHash)

			assert.ErrorIs(t, err, tc.expectedError)

			if tc.expectedError == nil {
				assert.FileExists(t, archivePath)
				removeArchive()
				assert.NoFileExists(t, archivePath)
			}
		})
	}
}
