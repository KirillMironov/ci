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

func (c Cloner) CloneRepository(url string) (sourceCodeArchivePath string, removeArchive func() error, err error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", nil, err
	}
	defer os.RemoveAll(dir)

	_, err = git.PlainClone(dir, false, &git.CloneOptions{URL: url})
	if err != nil {
		return "", nil, err
	}

	sourceCodeArchivePath, err = c.compress(dir)
	if err != nil {
		return "", nil, err
	}

	return sourceCodeArchivePath, func() error { return os.Remove(sourceCodeArchivePath) }, nil
}

func (Cloner) compress(srcPath string) (string, error) {
	archive, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	defer archive.Close()

	tw := tar.NewWriter(archive)

	err = filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsDir() {
			return nil
		}

		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, srcPath+string(os.PathSeparator))

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
		tw.Close()
		return "", err
	}

	return archive.Name(), tw.Close()
}
