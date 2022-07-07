package service

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestCloner_Clone(t *testing.T) {
	var cloner Cloner

	testCases := []struct {
		name   string
		url    string
		exists bool
	}{
		{
			name:   "repository",
			url:    "https://github.com/KirillMironov/kube",
			exists: true,
		},
		{
			name:   "not a repository",
			url:    "https://github.com",
			exists: false,
		},
		{
			name:   "not a repository",
			url:    "https://example.com",
			exists: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir, err := cloner.Clone(tc.url)
			assert.Equal(t, tc.exists, err == nil)

			files, err := ioutil.ReadDir(dir)
			assert.Equal(t, tc.exists, err == nil)
			assert.Equal(t, tc.exists, files != nil)
		})
	}
}
