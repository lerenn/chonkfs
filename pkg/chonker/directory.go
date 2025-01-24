package chonker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"slices"

	"github.com/lerenn/chonkfs/pkg/storage"
)

type DirectoryOption func(dir *directory)

func WithDirectoryLogger(logger *log.Logger) DirectoryOption {
	return func(dir *directory) {
		dir.logger = logger
	}
}

var _ Directory = (*directory)(nil)

type directory struct {
	storageDir storage.Directory
	opts       []DirectoryOption
	logger     *log.Logger
}

func NewDirectory(ctx context.Context, d storage.Directory, opts ...DirectoryOption) (*directory, error) {
	// Create a default directory
	dir := &directory{
		storageDir: d,
		opts:       opts,
		logger:     log.New(io.Discard, "", 0),
	}

	// Apply options
	for _, opt := range opts {
		opt(dir)
	}

	return dir, nil
}

func (dir *directory) GetAttributes(ctx context.Context) (DirectoryAttributes, error) {
	return DirectoryAttributes{}, nil
}

func (dir *directory) SetAttributes(ctx context.Context, attr DirectoryAttributes) error {
	return nil
}

func (dir *directory) checkIfFileOrDirectoryAlreadyExists(ctx context.Context, name string) error {
	// Check in directories
	_, err := dir.storageDir.GetDirectory(ctx, name)
	if err != nil && !errors.Is(err, storage.ErrDirectoryNotExists) {
		return fmt.Errorf("%w: %w", ErrChonker, err)
	} else if err == nil {
		return ErrAlreadyExists
	}

	// Check in files
	_, err = dir.storageDir.GetFile(ctx, name)
	if err != nil && !errors.Is(err, storage.ErrFileNotExists) {
		return fmt.Errorf("%w: %w", ErrChonker, err)
	} else if err == nil {
		return ErrAlreadyExists
	}

	return nil
}

func (dir *directory) CreateDirectory(ctx context.Context, name string) (Directory, error) {
	// Check if it doesn't not exist already
	if err := dir.checkIfFileOrDirectoryAlreadyExists(ctx, name); err != nil {
		return nil, err
	}

	// Create a new directory on storage
	nd, err := dir.storageDir.CreateDirectory(ctx, name)
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

func (dir *directory) GetDirectory(ctx context.Context, name string) (Directory, error) {
	// Check if this is not already a file
	_, err := dir.storageDir.GetFile(ctx, name)
	if err != nil && !errors.Is(err, storage.ErrFileNotExists) {
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	} else if err == nil {
		return nil, ErrNotDirectory
	}

	// Get and check if it exists
	d, err := dir.storageDir.GetDirectory(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrDirectoryNotExists) {
			return nil, ErrNoEntry
		}
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	return NewDirectory(ctx, d, dir.opts...)
}

func (dir *directory) GetFile(ctx context.Context, name string) (File, error) {
	// Get and check if it exists
	f, err := dir.storageDir.GetFile(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrFileNotExists) {
			return nil, ErrNoEntry
		}
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	// Get file info
	info, err := f.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	return NewFile(ctx, f, info.ChunkSize)
}

func (dir *directory) CreateFile(ctx context.Context, name string, chunkSize int) (File, error) {
	// Check if it doesn't not exist already
	if err := dir.checkIfFileOrDirectoryAlreadyExists(ctx, name); err != nil {
		return nil, err
	}

	// Create file on storage
	sf, err := dir.storageDir.CreateFile(ctx, name, chunkSize)
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

func (dir *directory) RemoveDirectory(ctx context.Context, name string) error {
	return dir.storageDir.RemoveDirectory(ctx, name)
}

func (dir *directory) RemoveFile(ctx context.Context, name string) error {
	return dir.storageDir.RemoveFile(ctx, name)
}

func (dir *directory) ListFiles(ctx context.Context) ([]string, error) {
	m, err := dir.storageDir.ListFiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	return slices.Collect(maps.Keys(m)), nil
}

func (dir *directory) RenameFile(ctx context.Context, name string, newParent Directory, newName string, noReplace bool) error {
	err := dir.storageDir.RenameFile(ctx, name, newParent.(*directory).storageDir, newName, noReplace)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, storage.ErrFileNotExists):
		return ErrNoEntry
	case errors.Is(err, storage.ErrFileAlreadyExists):
		return ErrAlreadyExists
	default:
		return fmt.Errorf("%w: %w", ErrChonker, err)
	}
}

func (dir *directory) ListDirectories(ctx context.Context) ([]string, error) {
	m, err := dir.storageDir.ListDirectories(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrChonker, err)
	}

	return slices.Collect(maps.Keys(m)), nil
}

func (dir *directory) RenameDirectory(ctx context.Context, name string, newParent Directory, newName string, noReplace bool) error {
	err := dir.storageDir.RenameDirectory(ctx, name, newParent.(*directory).storageDir, newName, noReplace)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, storage.ErrDirectoryNotExists):
		return ErrNoEntry
	case errors.Is(err, storage.ErrDirectoryAlreadyExists):
		return ErrAlreadyExists
	default:
		return fmt.Errorf("%w: %w", ErrChonker, err)
	}
}
