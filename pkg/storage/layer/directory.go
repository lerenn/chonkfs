package layer

import (
	"context"
	"errors"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.Directory = (*directory)(nil)

type DirectoryOptions struct {
	Underlayer storage.Directory
}

type directory struct {
	backend    storage.Directory
	underlayer storage.Directory
}

func NewDirectory(backend storage.Directory, opts *DirectoryOptions) *directory {
	d := &directory{
		backend: backend,
	}

	if opts != nil {
		d.underlayer = opts.Underlayer
	}

	return d
}

// CreateDirectory creates a directory.
func (d *directory) CreateDirectory(ctx context.Context, name string) (storage.Directory, error) {
	var child storage.Directory

	// Try to get from backend
	_, err := d.backend.GetDirectory(ctx, name)
	switch {
	case err == nil:
		// Directory already exists
		return nil, storage.ErrDirectoryAlreadyExists
	case errors.Is(err, storage.ErrIsFile):
		// File with the same name already exists
		return nil, storage.ErrFileAlreadyExists
	case !errors.Is(err, storage.ErrDirectoryNotFound):
		// Unexpected error
		return nil, err
	}

	// If there is an underlayer
	if d.underlayer != nil {
		// Create the directory on it
		child, err = d.underlayer.CreateDirectory(ctx, name)
		if err != nil {
			return nil, err
		}
	}

	// Create the directory on the backend
	dir, err := d.backend.CreateDirectory(ctx, name)
	if err != nil {
		return nil, err
	}

	// Return the new directory
	return NewDirectory(dir, &DirectoryOptions{
		Underlayer: child,
	}), nil
}

// GetInfo returns the directory info.
func (d *directory) GetInfo(_ context.Context) (info.Directory, error) {
	return info.Directory{}, nil
}

// ListFiles returns a map of files.
func (d *directory) ListFiles(ctx context.Context) (map[string]storage.File, error) {
	// Get local files
	backendFiles, err := d.backend.ListFiles(ctx)
	if err != nil {
		return nil, err
	}

	// Create the file representation
	files := make(map[string]storage.File, len(backendFiles))
	for n, f := range backendFiles {
		files[n] = newFile(f, 0, nil)
	}

	// Get underlayer files
	if d.underlayer != nil {
		underlayer, err := d.underlayer.ListFiles(ctx)
		if err != nil {
			return nil, err
		}

		// Merge the two maps
		for k, v := range underlayer {
			files[k] = v
		}
	}

	return files, nil
}

// GetDirectory returns a child directory.
func (d *directory) GetDirectory(ctx context.Context, name string) (storage.Directory, error) {
	var underlayer storage.Directory
	var err error

	// Get the directory from the underlayer
	if d.underlayer != nil {
		underlayer, err = d.underlayer.GetDirectory(ctx, name)
		if err != nil {
			return nil, err
		}
	}

	// Get the directory from the backend
	if _, err := d.backend.GetDirectory(ctx, name); err != nil {
		if !errors.Is(err, storage.ErrDirectoryNotFound) || underlayer == nil {
			return nil, err
		}
	}

	// Return the directory
	return NewDirectory(d.backend, &DirectoryOptions{
		Underlayer: underlayer,
	}), nil
}

// GetFile returns a child file.
func (d *directory) GetFile(ctx context.Context, name string) (storage.File, error) {
	var underlayer storage.File
	var err error

	// Get the directory from the underlayer
	if d.underlayer != nil {
		underlayer, err = d.underlayer.GetFile(ctx, name)
		if err != nil {
			return nil, err
		}
	}

	// Get the directory from the backend
	var info info.File
	backendFile, err := d.backend.GetFile(ctx, name)
	if err != nil {
		// If there is an error, or if the file doesn't exist and there is no underlayer, return error
		if !errors.Is(err, storage.ErrFileNotFound) || underlayer == nil {
			return nil, err
		}

		// Get the info from the underlayer
		info, err = underlayer.GetInfo(ctx)
		if err != nil {
			return nil, err
		}

		// Create a new file on the backend
		backendFile, err = d.backend.CreateFile(ctx, name, info.ChunkSize)
		if err != nil {
			return nil, err
		}
	} else {
		// Get the info from the backend
		info, err = backendFile.GetInfo(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Return the directory
	return newFile(backendFile, info.ChunkSize, &fileOptions{
		Underlayer: underlayer,
	}), nil
}

// ListDirectories returns a map of directories.
func (d *directory) ListDirectories(_ context.Context) (map[string]storage.Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateFile creates a file in the directory.
func (d *directory) CreateFile(ctx context.Context, name string, chunkSize int) (storage.File, error) {
	var child storage.File

	// Try to get from backend
	_, err := d.backend.GetFile(ctx, name)
	switch {
	case err == nil:
		// File already exists
		return nil, storage.ErrFileAlreadyExists
	case errors.Is(err, storage.ErrIsDirectory):
		// Directory with the same name already exists
		return nil, storage.ErrDirectoryAlreadyExists
	case !errors.Is(err, storage.ErrFileNotFound):
		// Unexpected error
		return nil, err
	}

	// If there is an underlayer
	if d.underlayer != nil {
		// Create the directory on it
		child, err = d.underlayer.CreateFile(ctx, name, chunkSize)
		if err != nil {
			return nil, err
		}
	}

	// Create the file on the backend
	backendFile, err := d.backend.CreateFile(ctx, name, chunkSize)
	if err != nil {
		return nil, err
	}

	// Return the new directory
	return newFile(backendFile, chunkSize, &fileOptions{
		Underlayer: child,
	}), nil
}

// RemoveDirectory removes a child directory of the directory.
func (d *directory) RemoveDirectory(ctx context.Context, name string) error {
	// Remove the directory from the underlayer
	if d.underlayer != nil {
		if err := d.underlayer.RemoveDirectory(ctx, name); err != nil {
			return err
		}
	}

	// Remove the directory from the backend
	err := d.backend.RemoveDirectory(ctx, name)
	if err == nil {
		return nil
	} else if errors.Is(err, storage.ErrDirectoryNotFound) && d.underlayer != nil {
		return nil
	}

	return err
}

// RemoveFile removes a child file of the directory.
func (d *directory) RemoveFile(_ context.Context, _ string) error {
	return fmt.Errorf("not implemented")
}

// RenameFile renames a child file of the directory.
func (d *directory) RenameFile(
	_ context.Context,
	_ string,
	_ storage.Directory,
	_ string,
	_ bool,
) error {
	return fmt.Errorf("not implemented")
}

// RenameDirectory renames a child directory of the directory.
func (d *directory) RenameDirectory(
	_ context.Context,
	_ string,
	_ storage.Directory,
	_ string,
	_ bool,
) error {
	return fmt.Errorf("not implemented")
}
