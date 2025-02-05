package files

import (
	"io"
	"os"
	"path/filepath"
)

// Local is an implementation of the Storage interface which works with the
// local disk on the current machine
type Local struct {
	maxFileSize int64
	basePath    string
}

// NewLocal creates a new Local filesystem with the given base path
// basePath is the base directory to save files to
// maxSize is the max number of bytes that a file can be
func NewLocal(basePath string, maxSize int64) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Local{basePath: p, maxFileSize: maxSize}, nil
}

// Save the contents of the Writer to the given path
// path is a relative path, basePath will be appended
func (l *Local) Save(path string, contents io.Reader) error {
	// creating the full path by joining the file path and base path
	fp := l.fullPath(path)

	// get the directory and make sure it exists
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return err
	}

	// if the file exists delete it
	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	// create the new file at the path
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	// copy the contents to the new file
	_, err = io.Copy(f, contents)
	if err != nil {
		return err
	}
	return nil
}

func (l *Local) Get(path string) (*os.File, error) {
	// get the full path for the file
	fp := l.fullPath(path)

	// open the file
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// returns the absolute path
func (l *Local) fullPath(path string) string {
	return filepath.Join(l.basePath, path)
}
