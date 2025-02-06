package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage/backend"
)

var _ Directory = (*directory)(nil)

type DirectoryOptions struct {
	Underlayer Directory
}

type directory struct {
	backend    backend.Directory
	underlayer Directory
}

func NewDirectory(backend backend.Directory, opts *DirectoryOptions) *directory {
	d := &directory{
		backend: backend,
	}

	if opts != nil {
		d.underlayer = opts.Underlayer
	}

	return d
}

// Underlayer returns the directory underlayer.
func (d *directory) Underlayer() Directory {
	return d.underlayer
}

// CreateDirectory creates a directory.
func (d *directory) CreateDirectory(ctx context.Context, name string) (Directory, error) {
	var child Directory

	// Try to get from backend
	_, err := d.backend.GetDirectory(ctx, name)
	switch {
	case err == nil:
		// Directory already exists
		return nil, ErrDirectoryAlreadyExists
	case errors.Is(err, backend.ErrIsFile):
		// File with the same name already exists
		return nil, ErrFileAlreadyExists
	case !errors.Is(err, backend.ErrNotFound):
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

// Info returns the directory info.
func (d *directory) Info(_ context.Context) (info.Directory, error) {
	return info.Directory{}, nil
}

// ListFiles returns a map of files.
func (d *directory) ListFiles(ctx context.Context) (map[string]File, error) {
	// Get local files
	backendFiles, err := d.backend.ListFiles(ctx)
	if err != nil {
		return nil, err
	}

	// Create the file representation
	files := make(map[string]File, len(backendFiles))
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
func (d *directory) GetDirectory(ctx context.Context, name string) (Directory, error) {
	var underlayer Directory
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
		if !errors.Is(err, backend.ErrNotFound) || underlayer == nil {
			return nil, err
		}
	}

	// Return the directory
	return NewDirectory(d.backend, &DirectoryOptions{
		Underlayer: underlayer,
	}), nil
}

// GetFile returns a child file.
func (d *directory) GetFile(ctx context.Context, name string) (File, error) {
	var underlayer File
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
		if !errors.Is(err, backend.ErrNotFound) || underlayer == nil {
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
func (d *directory) ListDirectories(_ context.Context) (map[string]Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateFile creates a file in the directory.
func (d *directory) CreateFile(ctx context.Context, name string, chunkSize int) (File, error) {
	var child File

	// Try to get from backend
	_, err := d.backend.GetFile(ctx, name)
	switch {
	case err == nil:
		// File already exists
		return nil, ErrFileAlreadyExists
	case errors.Is(err, backend.ErrIsDirectory):
		// Directory with the same name already exists
		return nil, ErrDirectoryAlreadyExists
	case !errors.Is(err, backend.ErrNotFound):
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
func (d *directory) RemoveDirectory(_ context.Context, _ string) error {
	return fmt.Errorf("not implemented")
}

// RemoveFile removes a child file of the directory.
func (d *directory) RemoveFile(_ context.Context, _ string) error {
	return fmt.Errorf("not implemented")
}

// RenameFile renames a child file of the directory.
func (d *directory) RenameFile(
	_ context.Context,
	_ string,
	_ Directory,
	_ string,
	_ bool,
) error {
	return fmt.Errorf("not implemented")
}

// RenameDirectory renames a child directory of the directory.
func (d *directory) RenameDirectory(
	_ context.Context,
	_ string,
	_ Directory,
	_ string,
	_ bool,
) error {
	return fmt.Errorf("not implemented")
}
