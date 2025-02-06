package mem

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/lerenn/chonkfs/pkg/storage/backend"
)

type Directory struct {
	directories map[string]*Directory
	files       map[string]*file
}

func NewDirectory() *Directory {
	return &Directory{
		directories: make(map[string]*Directory),
		files:       make(map[string]*file),
	}
}

func (d *Directory) CreateDirectory(_ context.Context, name string) error {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return fmt.Errorf("%w: %q", backend.ErrFileAlreadyExists, name)
	}

	// Check if a directory with this name already exists
	if _, ok := d.directories[name]; ok {
		return fmt.Errorf("%w: %q", backend.ErrDirectoryAlreadyExists, name)
	}

	d.directories[name] = NewDirectory()
	return nil
}

func (d *Directory) IsDirectory(_ context.Context, name string) error {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return fmt.Errorf("%w: %q", backend.ErrIsFile, name)
	}

	// Check if there is a directory with this name
	if _, ok := d.directories[name]; !ok {
		return fmt.Errorf("%w: %q", backend.ErrNotFound, name)
	}

	return nil
}

func (d *Directory) CreateFile(_ context.Context, name string, chunkSize int) error {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return fmt.Errorf("%w: %q", backend.ErrFileAlreadyExists, name)
	}

	// Check if there is a directory with this name
	if _, ok := d.directories[name]; ok {
		return fmt.Errorf("%w: %q", backend.ErrDirectoryAlreadyExists, name)
	}

	// Create the file if it does not exist
	if _, ok := d.files[name]; !ok {
		f, err := newFile(chunkSize)
		if err != nil {
			return err
		}
		d.files[name] = f
	}

	return nil
}

func (d *Directory) IsFile(_ context.Context, name string) error {
	// Check if there is a directory with this name
	if _, ok := d.directories[name]; ok {
		return fmt.Errorf("%w: %q", backend.ErrIsDirectory, name)
	}

	// Check if there is a file with this name
	if _, ok := d.files[name]; !ok {
		return fmt.Errorf("%w: %q", backend.ErrNotFound, name)
	}

	return nil
}

func (d *Directory) ListFiles(_ context.Context) ([]string, error) {
	return slices.Collect(maps.Keys(d.files)), nil
}
