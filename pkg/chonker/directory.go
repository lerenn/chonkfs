package chonker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"slices"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

type directoryOption func(dir *directory)

// WithDirectoryLogger is an option to set the logger of a directory.
//
//nolint:revive
func WithDirectoryLogger(logger *log.Logger) directoryOption {
	return func(dir *directory) {
		dir.logger = logger
	}
}

var _ Directory = (*directory)(nil)

type directory struct {
	storage storage.Directory
	opts    []directoryOption
	logger  *log.Logger
}

// NewDirectory creates a new directory.
func NewDirectory(_ context.Context, d storage.Directory, opts ...directoryOption) (Directory, error) {
	// Create a default directory
	dir := &directory{
		storage: d,
		opts:    opts,
		logger:  log.New(io.Discard, "", 0),
	}

	// Apply options
	for _, opt := range opts {
		opt(dir)
	}

	return dir, nil
}

// GetAttributes returns the attributes of the directory.
func (dir *directory) GetAttributes(_ context.Context) (DirectoryAttributes, error) {
	return DirectoryAttributes{}, nil
}

// SetAttributes sets the attributes of the directory.
func (dir *directory) SetAttributes(_ context.Context, _ DirectoryAttributes) error {
	return nil
}

func (dir *directory) checkIfFileOrDirectoryAlreadyExists(ctx context.Context, name string) error {
	// Check in directories
	_, err := dir.storage.GetDirectory(ctx, name)
	if err != nil && !errors.Is(err, storage.ErrDirectoryNotFound) {
		return fmt.Errorf("%w: %w", ErrChonker, err)
	} else if err == nil {
		return ErrAlreadyExists
	}

	// Check in files
	_, err = dir.storage.GetFile(ctx, name)
	if err != nil && !errors.Is(err, storage.ErrFileNotFound) {
		return fmt.Errorf("%w: %w", ErrChonker, err)
	} else if err == nil {
		return ErrAlreadyExists
	}

	return nil
}

// CreateDirectory creates a child directory to the directory.
func (dir *directory) CreateDirectory(ctx context.Context, name string) (Directory, error) {
	// Check if it doesn't not exist already
	if err := dir.checkIfFileOrDirectoryAlreadyExists(ctx, name); err != nil {
		return nil, err
	}

	// Create a new directory on storage
	nd, err := dir.storage.CreateDirectory(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	// Create a new directory
	d, err := NewDirectory(ctx, nd, dir.opts...)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// GetDirectory returns a child directory of the directory.
func (dir *directory) GetDirectory(ctx context.Context, name string) (Directory, error) {
	// Check if this is not already a file
	_, err := dir.storage.GetFile(ctx, name)
	if err != nil && !errors.Is(err, storage.ErrFileNotFound) && !errors.Is(err, storage.ErrIsDirectory) {
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	} else if err == nil {
		return nil, ErrNotDirectory
	}

	// Get and check if it exists
	d, err := dir.storage.GetDirectory(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrDirectoryNotFound) {
			return nil, ErrNoEntry
		}
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	return NewDirectory(ctx, d, dir.opts...)
}

// GetFile returns a child file of the directory.
func (dir *directory) GetFile(ctx context.Context, name string) (File, error) {
	// Get and check if it exists
	f, err := dir.storage.GetFile(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			return nil, ErrNoEntry
		}
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	// Get file info
	info, err := f.GetInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	return NewFile(ctx, f, info.ChunkSize)
}

// CreateFile creates a child file of the directory.
func (dir *directory) CreateFile(ctx context.Context, name string, chunkSize int) (File, error) {
	// Check if it doesn't not exist already
	if err := dir.checkIfFileOrDirectoryAlreadyExists(ctx, name); err != nil {
		return nil, err
	}

	// Create file on storage
	sf, err := dir.storage.CreateFile(ctx, name, info.File{
		ChunkSize: chunkSize,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	// Create file
	f, err := NewFile(ctx, sf, chunkSize,
		WithFileLogger(dir.logger))
	if err != nil {
		return nil, err
	}

	return f, nil
}

// RemoveDirectory removes a child directory of the directory.
func (dir *directory) RemoveDirectory(ctx context.Context, name string) error {
	return dir.storage.RemoveDirectory(ctx, name)
}

// RemoveFile removes a child file of the directory.
func (dir *directory) RemoveFile(ctx context.Context, name string) error {
	return dir.storage.RemoveFile(ctx, name)
}

// ListFiles returns the list of files in the directory.
func (dir *directory) ListFiles(ctx context.Context) ([]string, error) {
	m, err := dir.storage.ListFiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	return slices.Collect(maps.Keys(m)), nil
}

// RenameFile renames a child file of the directory.
func (dir *directory) RenameFile(
	ctx context.Context,
	name string,
	newParent Directory,
	newName string,
	noReplace bool,
) error {
	err := dir.storage.RenameFile(ctx, name, newParent.(*directory).storage, newName, noReplace)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, storage.ErrFileNotFound):
		return ErrNoEntry
	case errors.Is(err, storage.ErrFileAlreadyExists):
		return ErrAlreadyExists
	default:
		return fmt.Errorf("%w: %w", ErrChonker, err)
	}
}

// ListDirectories returns the list of directories in the directory.
func (dir *directory) ListDirectories(ctx context.Context) ([]string, error) {
	m, err := dir.storage.ListDirectories(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	return slices.Collect(maps.Keys(m)), nil
}

// RenameDirectory renames a child directory of the directory.
func (dir *directory) RenameDirectory(
	ctx context.Context,
	name string,
	newParent Directory,
	newName string,
	noReplace bool,
) error {
	err := dir.storage.RenameDirectory(ctx, name, newParent.(*directory).storage, newName, noReplace)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, storage.ErrDirectoryNotFound):
		return ErrNoEntry
	case errors.Is(err, storage.ErrDirectoryAlreadyExists):
		return ErrAlreadyExists
	default:
		return fmt.Errorf("%w: %w", ErrChonker, err)
	}
}
