package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/storage/backend"
)

var _ Directory = (*directory)(nil)

type DirectoryOptions struct {
	Underlayer Directory
}

type directory struct {
	backend    backend.BackEnd
	underlayer Directory
}

func NewDirectory(backend backend.BackEnd, opts *DirectoryOptions) *directory {
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
	err := d.backend.IsDirectory(ctx, name)
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
	if err := d.backend.CreateDirectory(ctx, name); err != nil {
		return nil, err
	}

	// Return the new directory
	return NewDirectory(d.backend, &DirectoryOptions{
		Underlayer: child,
	}), nil
}

// Info returns the directory info.
func (d *directory) Info(_ context.Context) (DirectoryInfo, error) {
	return DirectoryInfo{}, fmt.Errorf("not implemented")
}

// ListFiles returns a map of files.
func (d *directory) ListFiles(_ context.Context) (map[string]File, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetDirectory returns a child directory.
func (d *directory) GetDirectory(_ context.Context, _ string) (Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetFile returns a child file.
func (d *directory) GetFile(_ context.Context, _ string) (File, error) {
	return nil, fmt.Errorf("not implemented")
}

// ListDirectories returns a map of directories.
func (d *directory) ListDirectories(_ context.Context) (map[string]Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateFile creates a file in the directory.
func (d *directory) CreateFile(ctx context.Context, name string, chunkSize int) (File, error) {
	var child File

	// Try to get from backend
	err := d.backend.IsFile(ctx, name)
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
	if err := d.backend.CreateFile(ctx, name, chunkSize); err != nil {
		return nil, err
	}

	// Return the new directory
	return newFile(d.backend, chunkSize, &fileOptions{
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
