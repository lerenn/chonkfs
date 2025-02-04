package skeleton

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.Directory = (*Directory)(nil)

// DirectoryOptions represents the options that can be given to a Directory.
type DirectoryOptions struct {
	Underlayer storage.Directory
}

// Directory is a directory on disk.
type Directory struct {
	underlayer storage.Directory
}

// NewDirectory creates a new directory.
func NewDirectory(opts *DirectoryOptions) *Directory {
	d := &Directory{}

	if opts != nil {
		d.underlayer = opts.Underlayer
	}

	return d
}

// Underlayer returns the directory underlayer.
func (d *Directory) Underlayer() storage.Directory {
	return d.underlayer
}

// CreateDirectory creates a directory.
func (d *Directory) CreateDirectory(_ context.Context, _ string) (storage.Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

// Info returns the directory info.
func (d *Directory) Info(_ context.Context) (storage.DirectoryInfo, error) {
	return storage.DirectoryInfo{}, fmt.Errorf("not implemented")
}

// ListFiles returns a map of files.
func (d *Directory) ListFiles(_ context.Context) (map[string]storage.File, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetDirectory returns a child directory.
func (d *Directory) GetDirectory(_ context.Context, _ string) (storage.Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetFile returns a child file.
func (d *Directory) GetFile(_ context.Context, _ string) (storage.File, error) {
	return nil, fmt.Errorf("not implemented")
}

// ListDirectories returns a map of directories.
func (d *Directory) ListDirectories(_ context.Context) (map[string]storage.Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateFile creates a file in the directory.
func (d *Directory) CreateFile(_ context.Context, _ string, _ int) (storage.File, error) {
	return nil, fmt.Errorf("not implemented")
}

// RemoveDirectory removes a child directory of the directory.
func (d *Directory) RemoveDirectory(_ context.Context, _ string) error {
	return fmt.Errorf("not implemented")
}

// RemoveFile removes a child file of the directory.
func (d *Directory) RemoveFile(_ context.Context, _ string) error {
	return fmt.Errorf("not implemented")
}

// RenameFile renames a child file of the directory.
func (d *Directory) RenameFile(
	_ context.Context,
	_ string,
	_ storage.Directory,
	_ string,
	_ bool,
) error {
	return fmt.Errorf("not implemented")
}

// RenameDirectory renames a child directory of the directory.
func (d *Directory) RenameDirectory(
	_ context.Context,
	_ string,
	_ storage.Directory,
	_ string,
	_ bool,
) error {
	return fmt.Errorf("not implemented")
}
