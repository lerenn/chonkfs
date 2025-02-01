package skeleton

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.Directory = (*Directory)(nil)

// Directory is a directory on disk.
type Directory struct {
}

// NewDirectory creates a new directory.
func NewDirectory() *Directory {
	return &Directory{}
}

// Underlayer returns the directory underlayer.
func (d *Directory) Underlayer() storage.Directory {
	return nil
}

// CreateDirectory creates a directory.
func (d *Directory) CreateDirectory(_ context.Context, _ string) (storage.Directory, error) {
	return nil, nil
}

// Info returns the directory info.
func (d *Directory) Info(_ context.Context) (storage.DirectoryInfo, error) {
	return storage.DirectoryInfo{}, nil
}

// ListFiles returns a map of files.
func (d *Directory) ListFiles(_ context.Context) (map[string]storage.File, error) {
	return nil, nil
}

// GetDirectory returns a child directory.
func (d *Directory) GetDirectory(_ context.Context, _ string) (storage.Directory, error) {
	return nil, nil
}

// GetFile returns a child file.
func (d *Directory) GetFile(_ context.Context, _ string) (storage.File, error) {
	return nil, nil
}

// ListDirectories returns a map of directories.
func (d *Directory) ListDirectories(_ context.Context) (map[string]storage.Directory, error) {
	return nil, nil
}

// CreateFile creates a file in the directory.
func (d *Directory) CreateFile(_ context.Context, _ string, _ int) (storage.File, error) {
	return nil, nil
}

// RemoveDirectory removes a child directory of the directory.
func (d *Directory) RemoveDirectory(_ context.Context, _ string) error {
	return nil
}

// RemoveFile removes a child file of the directory.
func (d *Directory) RemoveFile(_ context.Context, _ string) error {
	return nil
}

func (d *Directory) checkIfFileOrDirectoryAlreadyExists(_ string) error {
	return nil
}

// RenameFile renames a child file of the directory.
func (d *Directory) RenameFile(
	_ context.Context,
	_ string,
	_ storage.Directory,
	_ string,
	_ bool,
) error {
	return nil
}

// RenameDirectory renames a child directory of the directory.
func (d *Directory) RenameDirectory(
	_ context.Context,
	_ string,
	_ storage.Directory,
	_ string,
	_ bool,
) error {
	return nil
}
