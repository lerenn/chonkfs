package skeleton

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.Directory = &directory{}

type directory struct{}

// NewDirectory creates a new directory representation.
func NewDirectory() storage.Directory {
	return &directory{}
}

// CreateDirectory creates a directory.
func (d *directory) CreateDirectory(_ context.Context, _ string) (storage.Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetDirectory returns a directory.
func (d *directory) GetDirectory(_ context.Context, _ string) (storage.Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetInfo returns the directory info.
func (d *directory) GetInfo(_ context.Context) (info.Directory, error) {
	return info.Directory{}, fmt.Errorf("not implemented")
}

// CreateFile creates a file.
func (d *directory) CreateFile(_ context.Context, _ string, info info.File) (storage.File, error) {
	_, _ = newFile(info)
	return nil, fmt.Errorf("not implemented")
}

// GetFile returns a file.
func (d *directory) GetFile(_ context.Context, _ string) (storage.File, error) {
	return nil, fmt.Errorf("not implemented")
}

// ListFiles returns a map of files.
func (d *directory) ListFiles(_ context.Context) (map[string]storage.File, error) {
	return nil, fmt.Errorf("not implemented")
}

// RemoveDirectory removes a directory.
func (d *directory) RemoveDirectory(_ context.Context, _ string) error {
	return fmt.Errorf("not implemented")
}

// ListDirectories returns a map of directories.
func (d *directory) ListDirectories(_ context.Context) (map[string]storage.Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

// RemoveFile removes a file.
func (d *directory) RemoveFile(_ context.Context, _ string) error {
	return fmt.Errorf("not implemented")
}

// RenameFile renames a file.
func (d *directory) RenameFile(
	_ context.Context,
	_ string,
	_ storage.Directory,
	_ string,
	_ bool,
) error {
	return fmt.Errorf("not implemented")
}

// RenameDirectory renames a directory.
func (d *directory) RenameDirectory(
	_ context.Context,
	_ string,
	_ storage.Directory,
	_ string,
	_ bool,
) error {
	return fmt.Errorf("not implemented")
}
