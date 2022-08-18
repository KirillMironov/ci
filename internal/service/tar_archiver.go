package service

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// TarArchiver used to work with tar archives.
type TarArchiver struct{}

// Compress compresses the given directory into a tar archive.
func (TarArchiver) Compress(dir string) (archivePath string, removeArchive func(), err error) {
	archive, err := os.CreateTemp("", "")
	if err != nil {
		return "", nil, err
	}
	defer archive.Close()

	var tw = tar.NewWriter(archive)

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, dir+string(os.PathSeparator))

		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tw, file)
		return err
	})
	if err != nil {
		tw.Close()
		return "", nil, err
	}

	return archive.Name(), func() { os.Remove(archive.Name()) }, tw.Close()
}
