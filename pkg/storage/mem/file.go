package mem

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

type file struct {
	data      [][]byte
	chunkSize int
}

func newFile(chunkSize int) (*file, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("%w: %d", storage.ErrInvalidChunkSize, chunkSize)
	}

	return &file{
		chunkSize: chunkSize,
	}, nil
}

func (f *file) GetInfo(_ context.Context) (info.File, error) {
	// Get last chunk size
	lastChunkSize := 0
	if len(f.data) > 0 {
		lastChunkSize = len(f.data[len(f.data)-1])
	}

	return info.File{
		Size:          (len(f.data)-1)*f.chunkSize + lastChunkSize,
		ChunkSize:     f.chunkSize,
		ChunksCount:   len(f.data),
		LastChunkSize: lastChunkSize,
	}, nil
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
