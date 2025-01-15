package chonkfs

import (
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
	dir := &Directory{
		backend: backend,
		logger:  log.Default(),
	}

	for _, o := range options {
		o(dir)
	}

	return dir
}
