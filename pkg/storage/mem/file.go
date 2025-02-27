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
	chunks        []*chunk
	chunkSize     int
	lastChunkSize int
}

func newFile(info info.File) (*file, error) {
	// Check chunk size
	if info.ChunkSize <= 0 {
		return nil, fmt.Errorf("%w: %d", storage.ErrInvalidChunkSize, info.ChunkSize)
	}

	// Create file representation
	f := &file{
		chunkSize:     info.ChunkSize,
		chunks:        make([]*chunk, info.ChunksCount),
		lastChunkSize: info.LastChunkSize,
	}

	return f, nil
}

func (f *file) GetInfo(_ context.Context) (info.File, error) {
	size := 0
	chunksCount := len(f.chunks)
	if chunksCount > 0 {
		size = (chunksCount-1)*f.chunkSize + f.lastChunkSize
	}

	return info.File{
		Size:          size,
		ChunkSize:     f.chunkSize,
		ChunksCount:   chunksCount,
		LastChunkSize: f.lastChunkSize,
	}, nil
}

func (f *file) checkReadWriteChunkParams(index int, offset int) error {
	// Check if chunk index is correct
	if index < 0 || index >= len(f.chunks) {
		return fmt.Errorf("%w: %d", storage.ErrInvalidChunkNb, index)
	}

	// Check if there is data to read
	if f.chunks[index] == nil {
		return fmt.Errorf("%w: %d", storage.ErrChunkNotFound, index)
	}

	// Check if offset is correct
	if offset < 0 || offset >= f.chunkSize {
		return fmt.Errorf("%w: %d", storage.ErrInvalidOffset, offset)
	}

	// Check if this is the last chunk, that the offset is correct$
	if index == len(f.chunks)-1 && offset >= f.chunks[index].Size {
		return fmt.Errorf("%w: %d", storage.ErrInvalidOffset, offset)
	}

	return nil
}

func (f *file) WriteChunk(_ context.Context, index int, data []byte, offset int) (int, error) {
	// Check params
	if err := f.checkReadWriteChunkParams(index, offset); err != nil {
		return 0, err
	}

	// Write data
	return copy(f.chunks[index].Data[offset:], data), nil
}

func (f *file) ReadChunk(_ context.Context, index int, data []byte, offset int) (int, error) {
	// Check params
	if err := f.checkReadWriteChunkParams(index, offset); err != nil {
		return 0, err
	}

	// Read data
	return copy(data, f.chunks[index].Data[offset:]), nil
}

func (f *file) ResizeChunksNb(_ context.Context, size int) error {
	// Check size is correct
	if size < 0 {
		return fmt.Errorf("%w: %d", storage.ErrInvalidChunkNb, size)
	}

	// Check the last chunk size is full
	if len(f.chunks) > 0 && f.lastChunkSize != f.chunkSize {
		return fmt.Errorf("%w", storage.ErrLastChunkNotFull)
	}

	// Resize chunks
	if size > len(f.chunks) {
		// Add chunks
		for i := len(f.chunks); i < size; i++ {
			f.chunks = append(f.chunks, &chunk{
				Data: make([]byte, f.chunkSize),
				Size: f.chunkSize,
			})
		}

		// Set last chunk size
		f.lastChunkSize = f.chunkSize
	} else {
		// Remove chunks
		f.chunks = f.chunks[:size]
	}

	return nil
}

func (f *file) ResizeLastChunk(_ context.Context, size int) (changed int, err error) {
	// Check size is correct
	if size < 0 || size > f.chunkSize {
		return 0, fmt.Errorf("%w: %d", storage.ErrInvalidChunkSize, size)
	}

	// Check if there is a last chunk
	if len(f.chunks) == 0 {
		return 0, fmt.Errorf("%w", storage.ErrNoChunk)
	}

	// Check if the last chunk is present
	lastChunk := f.chunks[len(f.chunks)-1]
	if lastChunk == nil {
		return 0, fmt.Errorf("%w", storage.ErrChunkNotFound)
	}

	// Resize last chunk
	oldSize := lastChunk.Size
	if size > oldSize {
		// Add data
		lastChunk.Data = append(lastChunk.Data, make([]byte, size-oldSize)...)
	} else {
		// Remove data
		lastChunk.Data = lastChunk.Data[:size]
	}

	// Set size
	lastChunk.Size = size
	f.lastChunkSize = size

	return size - oldSize, nil
}

func (f *file) ImportChunk(_ context.Context, index int, data []byte) error {
	// Check if chunk index is correct
	if index < 0 || index >= len(f.chunks) {
		return fmt.Errorf("%w: %d", storage.ErrInvalidChunkNb, index)
	}

	// Check if the chunk is empty
	if f.chunks[index] != nil {
		return fmt.Errorf("%w: %d", storage.ErrChunkAlreadyExists, index)
	}

	// Check if length of data is correct
	if (len(data) != f.chunkSize && index != len(f.chunks)-1) || len(data) > f.chunkSize {
		return fmt.Errorf("%w: %d", storage.ErrInvalidChunkSize, len(data))
	}

	// Import data
	f.chunks[index] = &chunk{
		Data: data,
		Size: len(data),
	}

	// If this is the last chunk, set the last chunk size
	if index == len(f.chunks)-1 {
		f.lastChunkSize = len(data)
	}

	return nil
}
