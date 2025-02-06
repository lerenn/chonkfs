package mem

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage/backend"
)

type file struct {
	data      [][]byte
	chunkSize int
}

func newFile(chunkSize int) (*file, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("%w: %d", backend.ErrInvalidChunkSize, chunkSize)
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
