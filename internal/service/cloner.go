package service

import (
	"archive/tar"
	"github.com/go-git/go-git/v5"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Cloner struct{}

func (c Cloner) CloneRepository(url string) (sourceCode io.ReadCloser, err error) {
	dir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	_, err = git.PlainClone(dir, false, &git.CloneOptions{URL: url})
	if err != nil {
		return nil, err
	}

	return c.compress(dir)
}

func (Cloner) compress(srcPath string) (io.ReadCloser, error) {
	archive, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return nil, err
	}

	tw := tar.NewWriter(archive)

	err = filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(filepath.ToSlash(path), srcPath+string(filepath.Separator))

		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tw, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return archive, tw.Close()
}
