package skeleton

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

type Directory struct{}

func NewDirectory() *Directory {
	return &Directory{}
}

func (d *Directory) CreateDirectory(_ context.Context, name string) (storage.Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *Directory) GetDirectory(_ context.Context, name string) (storage.Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *Directory) GetInfo(_ context.Context) (info.Directory, error) {
	return info.Directory{}, fmt.Errorf("not implemented")
}

func (d *Directory) CreateFile(_ context.Context, name string, chunkSize int) (storage.File, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *Directory) GetFile(ctx context.Context, name string) (storage.File, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *Directory) ListFiles(ctx context.Context) (map[string]storage.File, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *Directory) RemoveDirectory(ctx context.Context, name string) error {
	return fmt.Errorf("not implemented")
}

func (d *Directory) ListDirectories(ctx context.Context) (map[string]storage.Directory, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *Directory) RemoveFile(ctx context.Context, name string) error {
	return fmt.Errorf("not implemented")
}

func (d *Directory) RenameFile(ctx context.Context, name string, newParent storage.Directory, newName string, noReplace bool) error {
	return fmt.Errorf("not implemented")
}

func (d *Directory) RenameDirectory(ctx context.Context, name string, newParent storage.Directory, newName string, noReplace bool) error {
	return fmt.Errorf("not implemented")
}
