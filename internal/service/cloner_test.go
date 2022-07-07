package service

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestCloner_Clone(t *testing.T) {
	var cloner Cloner

	testCases := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "repository",
			url:         "https://github.com/KirillMironov/kube",
			expectError: true,
		},
		{
			name:        "not a repository",
			url:         "https://github.com",
			expectError: false,
		},
		{
			name:        "not a repository",
			url:         "https://example.com",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir, err := cloner.CloneRepository(tc.url)
			assert.Equal(t, tc.expectError, err == nil)

			files, err := ioutil.ReadDir(dir)
			assert.Equal(t, tc.expectError, err == nil)
			assert.Equal(t, tc.expectError, files != nil)
		})
	}
}
