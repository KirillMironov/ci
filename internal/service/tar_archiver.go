package service

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// TarArchiver used to work with tar archives.
type TarArchiver struct{}

// Compress compresses the given directory into a tar archive.
func (TarArchiver) Compress(dir string) (archivePath string, err error) {
	archive, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
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
		return "", err
	}

	return archive.Name(), tw.Close()
}

// FindFile finds a file in the given tar archive.
func (TarArchiver) FindFile(filename, archivePath string) ([]byte, error) {
	archive, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer archive.Close()

	var buf bytes.Buffer
	var tr = tar.NewReader(archive)

	for {
		header, err := tr.Next()
		if err != nil {
			return nil, err
		}
		if header.Name == filename {
			_, err = io.Copy(&buf, tr)
			if err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}
	}
}
