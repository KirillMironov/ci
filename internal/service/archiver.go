package service

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Archiver is a service that works with tar archives.
type Archiver struct{}

// Compress compresses the given directory into a tar archive.
func (Archiver) Compress(dir string) (archivePath string, remove func(), err error) {
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

// FindFile finds file by name in the given tar archive.
func (Archiver) FindFile(filename, archivePath string) ([]byte, error) {
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
