package mem

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

type chunk struct {
	Data []byte
	Size int
}

type file struct {
	chunks    []*chunk
	chunkSize int
}

func newFile(info info.File) (*file, error) {
	// Check chunk size
	if info.ChunkSize <= 0 {
		return nil, fmt.Errorf("%w: %d", storage.ErrInvalidChunkSize, info.ChunkSize)
	}

	// Create file representation
	f := &file{
		chunkSize: info.ChunkSize,
		chunks:    make([]*chunk, 0),
	}

	// Set chunks if required
	if info.ChunksCount > 0 {
		f.initChunks(info)
	}

	return f, nil
}

func (f *file) initChunks(info info.File) {
	chunks := make([]*chunk, 0)

	for i := 0; i < info.ChunksCount; i++ {
		var c *chunk
		if info.LastChunkSize > 0 && i == info.ChunksCount-1 {
			// Last chunk
			c = &chunk{
				Data: make([]byte, info.LastChunkSize),
				Size: info.LastChunkSize,
			}
			info.LastChunkSize = 0
		} else {
			// Regular chunk
			c = &chunk{
				Data: make([]byte, info.ChunkSize),
				Size: info.ChunkSize,
			}
		}
		chunks = append(chunks, c)
	}

	f.chunks = chunks
}

func (f *file) GetInfo(_ context.Context) (info.File, error) {
	// Get last chunk size
	lastChunkSize := 0
	if len(f.chunks) > 0 {
		lastChunkSize = f.chunks[len(f.chunks)-1].Size
	}

	return info.File{
		Size:          (len(f.chunks)-1)*f.chunkSize + lastChunkSize,
		ChunkSize:     f.chunkSize,
		ChunksCount:   len(f.chunks),
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
	if size < 0 {
		return fmt.Errorf("%w: %d", storage.ErrInvalidChunkNb, size)
	}

	if size > len(f.chunks) {
		// Add chunks
		for i := len(f.chunks); i < size; i++ {
			f.chunks = append(f.chunks, &chunk{
				Data: make([]byte, f.chunkSize),
				Size: f.chunkSize,
			})
		}
	} else {
		// Remove chunks
		f.chunks = f.chunks[:size]
	}

	return nil
}

func (f *file) ResizeLastChunk(ctx context.Context, size int) (changed int, err error) {
	// Check size is correct
	if size < 0 || size > f.chunkSize {
		return 0, fmt.Errorf("%w: %d", storage.ErrInvalidChunkSize, size)
	}

	// Check if there is a last chunk
	if len(f.chunks) == 0 {
		return 0, fmt.Errorf("%w", storage.ErrNoChunk)
	}

	// Get last chunk size
	lastChunkSize := 0
	lastChunk := f.chunks[len(f.chunks)-1]
	if len(f.chunks) > 0 {
		lastChunkSize = lastChunk.Size
	}

	// Resize last chunk
	if size > lastChunkSize {
		// Add data
		lastChunk.Data = append(lastChunk.Data, make([]byte, size-lastChunkSize)...)
	} else {
		// Remove data
		lastChunk.Data = lastChunk.Data[:size]
	}
	lastChunk.Size = size

	return size - lastChunkSize, nil
}
