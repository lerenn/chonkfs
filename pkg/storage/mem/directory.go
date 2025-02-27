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

// NewDirectory creates a new directory.
func NewDirectory() storage.Directory {
	return &directory{
		directories: make(map[string]storage.Directory),
		files:       make(map[string]storage.File),
	}
}

// CreateDirectory creates a directory.
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

// GetDirectory returns a directory.
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

// GetInfo returns the directory info.
func (d *directory) GetInfo(_ context.Context) (info.Directory, error) {
	return info.Directory{}, nil
}

// CreateFile creates a file.
func (d *directory) CreateFile(_ context.Context, name string, info info.File) (storage.File, error) {
	// Check if there is a file with this name
	if _, ok := d.files[name]; ok {
		return nil, fmt.Errorf("%w: %q", storage.ErrFileAlreadyExists, name)
	}

	// Check if there is a directory with this name
	if _, ok := d.directories[name]; ok {
		return nil, fmt.Errorf("%w: %q", storage.ErrDirectoryAlreadyExists, name)
	}

	// Create the file
	f, err := newFile(info)
	if err != nil {
		return nil, err
	}

	// Store the file
	d.files[name] = f

	return f, nil
}

// GetFile returns a file.
func (d *directory) GetFile(_ context.Context, name string) (storage.File, error) {
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

// ListFiles returns a map of files.
func (d *directory) ListFiles(_ context.Context) (map[string]storage.File, error) {
	return maps.Clone(d.files), nil
}

// RemoveDirectory removes a directory.
func (d *directory) RemoveDirectory(_ context.Context, name string) error {
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

// ListDirectories returns a map of directories.
func (d *directory) ListDirectories(_ context.Context) (map[string]storage.Directory, error) {
	return maps.Clone(d.directories), nil
}

// RemoveFile removes a file.
func (d *directory) RemoveFile(_ context.Context, name string) error {
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

// RenameFile renames a file.
func (d *directory) RenameFile(
	_ context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// Check if there is a file with this name
	if _, ok := d.files[name]; !ok {
		return fmt.Errorf("%w: %q", storage.ErrFileNotFound, name)
	}

	// Check if there is a directory with the new name
	_, directoryExists := newParent.(*directory).directories[newName]
	if directoryExists {
		if noReplace {
			return fmt.Errorf("%w: %q", storage.ErrDirectoryAlreadyExists, newName)
		}

		// Delete directory
		delete(newParent.(*directory).directories, newName)
	}

	// Check if there is a file with the new name
	_, fileExists := newParent.(*directory).files[newName]
	if fileExists {
		if noReplace {
			return fmt.Errorf("%w: %q", storage.ErrFileAlreadyExists, newName)
		}

		// Delete file
		delete(newParent.(*directory).files, newName)
	}

	// Move the file
	newParent.(*directory).files[newName] = d.files[name]
	delete(d.files, name)

	return nil
}

// RenameDirectory renames a directory.
func (d *directory) RenameDirectory(
	_ context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// Check if there is a directory with this name
	if _, ok := d.directories[name]; !ok {
		return fmt.Errorf("%w: %q", storage.ErrDirectoryNotFound, name)
	}

	// Check if there is a directory with the new name
	_, directoryExists := newParent.(*directory).directories[newName]
	if directoryExists {
		if noReplace {
			return fmt.Errorf("%w: %q", storage.ErrDirectoryAlreadyExists, newName)
		}

		// Delete directory
		delete(newParent.(*directory).directories, newName)
	}

	// Check if there is a file with the new name
	_, fileExists := newParent.(*directory).files[newName]
	if fileExists {
		if noReplace {
			return fmt.Errorf("%w: %q", storage.ErrFileAlreadyExists, newName)
		}

		// Delete file
		delete(newParent.(*directory).files, newName)
	}

	// Move the directory
	newParent.(*directory).directories[newName] = d.directories[name]
	delete(d.directories, name)

	return nil
}
