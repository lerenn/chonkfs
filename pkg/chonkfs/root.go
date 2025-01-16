package chonkfs

import (
	"io"
	"log"

	"github.com/lerenn/chonkfs/pkg/backends"
)

type DirectoryOption func(dir *Directory)

func WithLogger(logger *log.Logger) DirectoryOption {
	return func(dir *Directory) {
		dir.logger = logger
	}
}

func New(backend backends.Directory, options ...DirectoryOption) *Directory {
	// Create a new directory
	dir := &Directory{
		backend: backend,
		logger:  log.New(io.Discard, "", 0),
	}

	// Apply options
	for _, o := range options {
		o(dir)
	}

	return dir
}
