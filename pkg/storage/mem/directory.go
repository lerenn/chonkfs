package mem

import (
	"context"
	"fmt"
	"maps"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

type directory struct {
	directories map[string]storage.Directory
	files       map[string]storage.File
}

func NewDirectory() *directory {
	return &directory{
		directories: make(map[string]storage.Directory),
		files:       make(map[string]storage.File),
	}
}

func (d *directory) CreateDirectory(_ context.Context, name string) (storage.Directory, error) {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return nil, fmt.Errorf("%w: %q", storage.ErrFileAlreadyExists, name)
	}

	// Check if a directory with this name already exists
	if _, ok := d.directories[name]; ok {
		return nil, fmt.Errorf("%w: %q", storage.ErrDirectoryAlreadyExists, name)
	}

	// Create directory and store it
	nd := NewDirectory()
	d.directories[name] = nd

	return nd, nil
}

func (d *directory) GetDirectory(_ context.Context, name string) (storage.Directory, error) {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return nil, fmt.Errorf("%w: %q", storage.ErrIsFile, name)
	}

	// Check if there is a directory with this name
	nd, ok := d.directories[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", storage.ErrDirectoryNotFound, name)
	}

	return nd, nil
}

func (d *directory) GetInfo(_ context.Context) (info.Directory, error) {
	return info.Directory{}, nil
}

func (d *directory) CreateFile(_ context.Context, name string, chunkSize int) (storage.File, error) {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return nil, fmt.Errorf("%w: %q", storage.ErrFileAlreadyExists, name)
	}

	// Check if there is a directory with this name
	if _, ok := d.directories[name]; ok {
		return nil, fmt.Errorf("%w: %q", storage.ErrDirectoryAlreadyExists, name)
	}

	// Create the file
	f, err := newFile(chunkSize)
	if err != nil {
		return nil, err
	}

	d.files[name] = f

	return f, nil
}

func (d *directory) GetFile(ctx context.Context, name string) (storage.File, error) {
	// Check if there is a directory with this name
	if _, ok := d.directories[name]; ok {
		return nil, fmt.Errorf("%w: %q", storage.ErrIsDirectory, name)
	}

	// Check if there is a file with this name
	f, ok := d.files[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", storage.ErrFileNotFound, name)
	}

	return f, nil
}

func (d *directory) ListFiles(ctx context.Context) (map[string]storage.File, error) {
	return maps.Clone(d.files), nil
}

func (d *directory) RemoveDirectory(ctx context.Context, name string) error {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return fmt.Errorf("%w: %q", storage.ErrIsFile, name)
	}

	// Check if there is a directory with this name
	if _, ok := d.directories[name]; !ok {
		return fmt.Errorf("%w: %q", storage.ErrDirectoryNotFound, name)
	}

	// Remove the directory
	delete(d.directories, name)

	return nil
}

func (d *directory) ListDirectories(ctx context.Context) (map[string]storage.Directory, error) {
	return maps.Clone(d.directories), nil
}

func (d *directory) RemoveFile(ctx context.Context, name string) error {
	// Check if there is a directory with this name
	if _, ok := d.directories[name]; ok {
		return fmt.Errorf("%w: %q", storage.ErrIsDirectory, name)
	}

	// Check if there is a file with this name
	if _, ok := d.files[name]; !ok {
		return fmt.Errorf("%w: %q", storage.ErrFileNotFound, name)
	}

	// Remove the file
	delete(d.files, name)

	return nil
}

func (d *directory) RenameFile(ctx context.Context, name string, newParent storage.Directory, newName string, noReplace bool) error {
	// Check if there is a file with this name
	if _, ok := d.files[name]; !ok {
		return fmt.Errorf("%w: %q", storage.ErrFileNotFound, name)
	}

	// Check if there is a directory with the new name
	_, directoryExists := newParent.(*directory).directories[newName]
	if directoryExists {
		return fmt.Errorf("%w: %q", storage.ErrDirectoryAlreadyExists, newName)
	}

	// Check if there is a file with the new name
	_, fileExists := newParent.(*directory).files[newName]
	if fileExists && !noReplace {
		return fmt.Errorf("%w: %q", storage.ErrFileAlreadyExists, newName)
	}

	// Move the file
	newParent.(*directory).files[newName] = d.files[name]
	delete(d.files, name)

	return nil
}

func (d *directory) RenameDirectory(ctx context.Context, name string, newParent storage.Directory, newName string, noReplace bool) error {
	return fmt.Errorf("not implemented")
}
