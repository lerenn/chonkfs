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

func (f *file) ImportChunk(_ context.Context, _ int, _ []byte) error {
	return fmt.Errorf("not implemented")
}

func (f *file) GetInfo(_ context.Context) (info.File, error) {
	return info.File{}, fmt.Errorf("not implemented")
}

func (f *file) WriteChunk(_ context.Context, _ int, _ []byte, _ int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (f *file) ReadChunk(_ context.Context, _ int, _ []byte, _ int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (f *file) ResizeChunksNb(_ context.Context, _ int) error {
	return fmt.Errorf("not implemented")
}

func (f *file) ResizeLastChunk(_ context.Context, _ int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}
