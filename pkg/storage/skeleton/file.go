package skeleton

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/info"
)

type file struct{}

func newFile() (*file, error) {
	return &file{}, nil
}

func (f *file) GetInfo(_ context.Context) (info.File, error) {
	return info.File{}, fmt.Errorf("not implemented")
}

func (f *file) WriteChunk(ctx context.Context, chunkIndex int, data []byte, offset int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (f *file) ReadChunk(ctx context.Context, chunkIndex int, data []byte, offset int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (f *file) ResizeChunksNb(ctx context.Context, size int) error {
	return fmt.Errorf("not implemented")
}

func (f *file) ResizeLastChunk(ctx context.Context, size int) (changed int, err error) {
	return 0, fmt.Errorf("not implemented")
}
