package mem

import (
	"context"
	"errors"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.Directory = (*Directory)(nil)

// Directory is a directory in memory.
type Directory struct {
	directories map[string]*Directory
	files       map[string]*File
}

// NewDirectory creates a new directory.
func NewDirectory() *Directory {
	return &Directory{
		directories: make(map[string]*Directory),
		files:       make(map[string]*File),
	}
}

// CreateDirectory creates a directory.
func (d *Directory) CreateDirectory(_ context.Context, name string) (storage.Directory, error) {
	if _, exist := d.directories[name]; exist {
		return nil, storage.ErrDirectoryAlreadyExists
	}

	nd := NewDirectory()
	d.directories[name] = nd
	return nd, nil
}

// Info returns the directory info.
func (d *Directory) Info(_ context.Context) (storage.DirectoryInfo, error) {
	return storage.DirectoryInfo{}, nil
}

// ListFiles returns a map of files.
func (d *Directory) ListFiles(_ context.Context) (map[string]storage.File, error) {
	m := make(map[string]storage.File, len(d.files))
	for p, f := range d.files {
		m[p] = f
	}
	return m, nil
}

// GetDirectory returns a child directory.
func (d *Directory) GetDirectory(_ context.Context, name string) (storage.Directory, error) {
	dir, ok := d.directories[name]
	if !ok {
		return nil, storage.ErrDirectoryNotExists
	}
	return dir, nil
}

// GetFile returns a child file.
func (d *Directory) GetFile(_ context.Context, name string) (storage.File, error) {
	f, ok := d.files[name]
	if !ok {
		return nil, storage.ErrFileNotExists
	}
	return f, nil
}

// ListDirectories returns a map of directories.
func (d *Directory) ListDirectories(_ context.Context) (map[string]storage.Directory, error) {
	m := make(map[string]storage.Directory, len(d.directories))
	for p, dir := range d.directories {
		m[p] = dir
	}
	return m, nil
}

// CreateFile creates a file in the directory.
func (d *Directory) CreateFile(_ context.Context, name string, chunkSize int) (storage.File, error) {
	if _, exist := d.files[name]; exist {
		return nil, fmt.Errorf("couldn't create file %q: %w", name, storage.ErrFileAlreadyExists)
	}

	f := newFile(chunkSize)
	d.files[name] = f
	return f, nil
}

// RemoveDirectory removes a child directory of the directory.
func (d *Directory) RemoveDirectory(_ context.Context, name string) error {
	if _, exist := d.directories[name]; !exist {
		return storage.ErrDirectoryNotExists
	}
	delete(d.directories, name)
	return nil
}

// RemoveFile removes a child file of the directory.
func (d *Directory) RemoveFile(_ context.Context, name string) error {
	if _, exist := d.files[name]; !exist {
		return storage.ErrFileNotExists
	}
	delete(d.files, name)
	return nil
}

func (d *Directory) checkIfFileOrDirectoryAlreadyExists(name string) error {
	if _, exist := d.directories[name]; exist {
		return fmt.Errorf("%w: %q", storage.ErrDirectoryAlreadyExists, name)
	}
	if _, exist := d.files[name]; exist {
		return fmt.Errorf("%w: %q", storage.ErrFileAlreadyExists, name)
	}
	return nil
}

// RenameFile renames a child file of the directory.
func (d *Directory) RenameFile(
	_ context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// Get the directory or the file
	f, fileExist := d.files[name]
	if !fileExist {
		return storage.ErrFileNotExists
	}

	// Check if it doesn't not exist already
	if err := newParent.(*Directory).checkIfFileOrDirectoryAlreadyExists(newName); err != nil {
		if errors.Is(err, storage.ErrFileAlreadyExists) && noReplace {
			// If noReplace is set and the file already exists, return error
			return err
		} else if !errors.Is(err, storage.ErrFileAlreadyExists) {
			// If another error, then return
			return err
		}
	}

	// Add it to new parent and remove it from current parent
	newParent.(*Directory).files[newName] = f
	delete(d.files, name)

	return nil
}

// RenameDirectory renames a child directory of the directory.
func (d *Directory) RenameDirectory(
	_ context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// Get the directory or the file
	dir, dirExist := d.directories[name]
	if !dirExist {
		return storage.ErrFileNotExists
	}

	// Check if it doesn't not exist already
	if err := newParent.(*Directory).checkIfFileOrDirectoryAlreadyExists(newName); err != nil {
		if errors.Is(err, storage.ErrDirectoryAlreadyExists) && noReplace {
			// If noReplace is set and the file already exists, return error
			return err
		} else if !errors.Is(err, storage.ErrDirectoryAlreadyExists) {
			// If another error, then return
			return err
		}
	}

	// Add it to new parent and remove it from current parent
	newParent.(*Directory).directories[newName] = dir
	delete(d.directories, name)

	return nil
}
