package mem

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/storage"
)

type File struct {
	data      [][]byte
	chunkSize int
}

func newFile(chunkSize int) *File {
	return &File{
		data:      make([][]byte, 0),
		chunkSize: chunkSize,
	}
}

func (f *File) Info(ctx context.Context) (storage.FileInfo, error) {
	return storage.FileInfo{
		ChunkSize: f.chunkSize,
	}, nil
}

func (f *File) ReadChunk(ctx context.Context, chunkIndex int, data []byte, start int, end *int) (int, error) {
	// Check if the chunk index is valid
	if chunkIndex < 0 || chunkIndex >= len(f.data) {
		return 0, storage.ErrInvalidChunkNb
	}

	// Check if the start is valid
	if start < 0 || start >= len(f.data[chunkIndex]) {
		return 0, fmt.Errorf("%w: start is %d", storage.ErrInvalidStartOffset, start)
	}

	// Check if the end is valid
	if end != nil && (*end < 0 || *end > len(f.data[chunkIndex])) {
		return 0, fmt.Errorf("%w: end is %d", storage.ErrInvalidEndOffset, start)
	}

	// Set the end if it is nil
	if end == nil {
		end = new(int)
		*end = len(f.data[chunkIndex])
	}

	// Read the data
	return copy(data, f.data[chunkIndex][start:*end]), nil
}

func (f *File) ChunksCount(ctx context.Context) (int, error) {
	return len(f.data), nil
}

func (f *File) WriteChunk(ctx context.Context, chunkIndex int, start int, end *int, data []byte) (int, error) {

	// Check if the chunk index is valid
	if chunkIndex < 0 || chunkIndex >= len(f.data) {
		return 0, storage.ErrInvalidChunkNb
	}

	// Check if the start is valid
	if start < 0 || start >= len(f.data[chunkIndex]) {
		return 0, fmt.Errorf("%w: start is %d", storage.ErrInvalidStartOffset, start)
	}

	// Check if the end is valid
	if end != nil && (*end < 0 || *end > len(f.data[chunkIndex])) {
		return 0, fmt.Errorf("%w: end is %d", storage.ErrInvalidEndOffset, start)
	}

	// Set the end if it is nil
	if end == nil {
		end = new(int)
		*end = len(f.data[chunkIndex])
	}

	// Write the data
	return copy(f.data[chunkIndex][start:*end], data), nil
}

func (f *File) ResizeChunksNb(ctx context.Context, nb int) error {
	// Check if the number of chunks is valid
	if nb < 0 {
		return storage.ErrInvalidChunkNb
	}

	if nb > len(f.data) {
		// Check if the last chunk is full
		if len(f.data) > 0 && len(f.data[len(f.data)-1]) != f.chunkSize {
			return storage.ErrLastChunkNotFull
		}

		// Add chunks
		for i := len(f.data); i < nb; i++ {
			f.data = append(f.data, make([]byte, f.chunkSize))
		}
	} else if nb < len(f.data) {
		// Remove chunks
		f.data = f.data[:nb]
	}

	return nil
}

func (f *File) ResizeLastChunk(ctx context.Context, size int) (changed int, err error) {
	// Check if the size is valid
	if size < 0 || size > f.chunkSize {
		return 0, storage.ErrInvalidChunkSize
	}

	// Check if there is a chunk to resize
	if len(f.data) == 0 {
		return 0, storage.ErrNoChunk
	}

	lastChunkSize := len(f.data[len(f.data)-1])
	if size < lastChunkSize {
		// Truncate the last chunk
		f.data[len(f.data)-1] = f.data[len(f.data)-1][:size]
		return -size, nil
	} else if size > lastChunkSize {
		// Add data to the last chunk
		toAdd := size - lastChunkSize
		f.data[len(f.data)-1] = append(f.data[len(f.data)-1], make([]byte, toAdd)...)
		return toAdd, nil
	}

	return 0, nil
}

func (f *File) Size(ctx context.Context) (int, error) {
	// Check if there is no data
	size := len(f.data)
	if size == 0 {
		return 0, nil
	}

	// Return the count of all chunks except the last one, multiplied by the chunk size
	// + the length of the last chunk
	return (size-1)*f.chunkSize + len(f.data[size-1]), nil
}

func (f *File) LastChunkSize(ctx context.Context) (int, error) {
	// Check if there is no data
	if len(f.data) == 0 {
		return 0, storage.ErrNoChunk
	}

	return len(f.data[len(f.data)-1]), nil
}
