package skeleton

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.File = &file{}

type file struct{}

func newFile(_ info.File) (*file, error) {
	return &file{}, nil
}

func (f *file) ImportChunk(ctx context.Context, index int, data []byte) error {
	return fmt.Errorf("not implemented")
}

func (f *file) GetInfo(_ context.Context) (info.File, error) {
	return info.File{}, fmt.Errorf("not implemented")
}

func (f *file) WriteChunk(ctx context.Context, index int, data []byte, offset int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (f *file) ReadChunk(ctx context.Context, index int, data []byte, offset int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (f *file) ResizeChunksNb(ctx context.Context, size int) error {
	return fmt.Errorf("not implemented")
}

func (f *file) ResizeLastChunk(ctx context.Context, size int) (changed int, err error) {
	return 0, fmt.Errorf("not implemented")
}
