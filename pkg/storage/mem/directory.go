package mem

import (
	"context"
	"errors"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.Directory = (*Directory)(nil)

type DirectoryOptions struct {
	Underlayer storage.Directory
}

// Directory is a directory in memory.
type Directory struct {
	directories map[string]*Directory
	files       map[string]*file
	opts        *DirectoryOptions
}

// NewDirectory creates a new directory.
func NewDirectory(opts *DirectoryOptions) *Directory {
	return &Directory{
		directories: make(map[string]*Directory),
		files:       make(map[string]*file),
		opts:        opts,
	}
}

func (d *Directory) Underlayer() storage.Directory {
	if d.opts == nil {
		return nil
	}

	return d.opts.Underlayer
}

// CreateDirectory creates a directory.
func (d *Directory) CreateDirectory(ctx context.Context, name string) (storage.Directory, error) {
	var childUnderlayer storage.Directory
	var err error

	// Check if the directory already exists
	if u := d.Underlayer(); u != nil {
		// If there is an underlayer, then creates it here: it will check if
		// the directory already exists
		childUnderlayer, err = u.CreateDirectory(ctx, name)
		if err != nil {
			return nil, err
		}
	} else if _, exist := d.directories[name]; exist {
		// If already exists, then return an error
		return nil, storage.ErrDirectoryAlreadyExists
	}

	// Create the new directory with its underlayer
	nd := NewDirectory(&DirectoryOptions{
		Underlayer: childUnderlayer,
	})
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
func (d *Directory) CreateFile(ctx context.Context, name string, chunkSize int) (storage.File, error) {
	var childUnderlayer storage.File
	var err error

	// Check if the file already exists
	if u := d.Underlayer(); u != nil {
		// Create the file in the underlayer, that will check that the file already exists
		childUnderlayer, err = u.CreateFile(ctx, name, chunkSize)
		if err != nil {
			return nil, err
		}
	} else if _, exist := d.files[name]; exist {
		// If the file already exists, return an error
		return nil, fmt.Errorf("couldn't create file %q: %w", name, storage.ErrFileAlreadyExists)
	}

	f := newFile(chunkSize, &fileOptions{
		Underlayer: childUnderlayer,
	})
	d.files[name] = f
	return f, nil
}

// RemoveDirectory removes a child directory of the directory.
func (d *Directory) RemoveDirectory(ctx context.Context, name string) error {
	// Check if the directory already exists
	if u := d.Underlayer(); u != nil {
		// If there is an underlayer, then creates it here: it will check if
		// the directory already exists
		if err := u.RemoveDirectory(ctx, name); err != nil {
			return err
		}
	} else if _, exist := d.directories[name]; !exist {
		// If there is no corresponding directory, then return an error
		return storage.ErrDirectoryNotExists
	}

	// Actually delete the directory
	delete(d.directories, name)

	return nil
}

// RemoveFile removes a child file of the directory.
func (d *Directory) RemoveFile(ctx context.Context, name string) error {
	// Check if the file exists
	if u := d.Underlayer(); u != nil {
		// If there is an underlayer, remove it there, it will check if the file exists
		if err := u.RemoveFile(ctx, name); err != nil {
			return err
		}
	} else if _, exist := d.files[name]; !exist {
		// If the file doesn't exists, return an error
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
	ctx context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// If there is an underlayer, then rename the file here first
	if u := d.Underlayer(); u != nil {
		if err := u.RenameFile(ctx, name, newParent.Underlayer(), newName, noReplace); err != nil {
			return err
		}
	}

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
	ctx context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// If there is an underlayer, then rename the directory here first
	if u := d.Underlayer(); u != nil {
		if err := u.RenameDirectory(ctx, name, newParent.Underlayer(), newName, noReplace); err != nil {
			return err
		}
	}

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
