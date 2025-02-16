package files

import (
	"io"
	"os"
)

// Storage defines the behaviour for file operations
// Implementation may be of using local file store or using cloud storage such as AWS S3
type Storage interface {
	Save(path string, r io.Reader) error
	Get(path string) (*os.File, error)
}
