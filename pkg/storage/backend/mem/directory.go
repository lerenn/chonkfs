package mem

import (
	"context"
	"fmt"
	"maps"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage/backend"
)

type Directory struct {
	directories map[string]backend.Directory
	files       map[string]backend.File
}

func NewDirectory() *Directory {
	return &Directory{
		directories: make(map[string]backend.Directory),
		files:       make(map[string]backend.File),
	}
}

func (d *Directory) CreateDirectory(_ context.Context, name string) (backend.Directory, error) {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return nil, fmt.Errorf("%w: %q", backend.ErrFileAlreadyExists, name)
	}

	// Check if a directory with this name already exists
	if _, ok := d.directories[name]; ok {
		return nil, fmt.Errorf("%w: %q", backend.ErrDirectoryAlreadyExists, name)
	}

	// Create directory and store it
	nd := NewDirectory()
	d.directories[name] = nd

	return nd, nil
}

func (d *Directory) GetDirectory(_ context.Context, name string) (backend.Directory, error) {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return nil, fmt.Errorf("%w: %q", backend.ErrIsFile, name)
	}

	// Check if there is a directory with this name
	nd, ok := d.directories[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", backend.ErrNotFound, name)
	}

	return nd, nil
}

func (d *Directory) GetInfo(_ context.Context) (info.Directory, error) {
	return info.Directory{}, nil
}

func (d *Directory) CreateFile(_ context.Context, name string, chunkSize int) (backend.File, error) {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return nil, fmt.Errorf("%w: %q", backend.ErrFileAlreadyExists, name)
	}

	// Check if there is a directory with this name
	if _, ok := d.directories[name]; ok {
		return nil, fmt.Errorf("%w: %q", backend.ErrDirectoryAlreadyExists, name)
	}

	// Create the file
	f, err := newFile(chunkSize)
	if err != nil {
		return nil, err
	}

	d.files[name] = f

	return f, nil
}

func (d *Directory) GetFile(ctx context.Context, name string) (backend.File, error) {
	// Check if there is a directory with this name
	if _, ok := d.directories[name]; ok {
		return nil, fmt.Errorf("%w: %q", backend.ErrIsDirectory, name)
	}

	// Check if there is a file with this name
	f, ok := d.files[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", backend.ErrNotFound, name)
	}

	return f, nil
}

func (d *Directory) ListFiles(ctx context.Context) (map[string]backend.File, error) {
	return maps.Clone(d.files), nil
}

func (d *Directory) RemoveDirectory(ctx context.Context, name string) error {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return fmt.Errorf("%w: %q", backend.ErrIsFile, name)
	}

	// Check if there is a directory with this name
	if _, ok := d.directories[name]; !ok {
		return fmt.Errorf("%w: %q", backend.ErrNotFound, name)
	}

	// Remove the directory
	delete(d.directories, name)

	return nil
}
