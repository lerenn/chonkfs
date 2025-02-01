package mem

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/storage"
)

type FileOptions struct {
	Underlayer storage.File
}

// File is a file in chunks in memory.
type File struct {
	data      [][]byte
	chunkSize int
	opts      *FileOptions
}

func newFile(chunkSize int, opts *FileOptions) *File {
	return &File{
		data:      make([][]byte, 0),
		chunkSize: chunkSize,
		opts:      opts,
	}
}

func (f *File) Underlayer() storage.File {
	if f.opts == nil {
		return nil
	}

	return f.opts.Underlayer
}

// Info returns the file info.
func (f *File) Info(_ context.Context) (storage.FileInfo, error) {
	return storage.FileInfo{
		ChunkSize: f.chunkSize,
	}, nil
}

// ReadChunk reads data from a chunk.
func (f *File) ReadChunk(_ context.Context, chunkIndex int, data []byte, start int, end *int) (int, error) {
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

// ChunksCount returns the number of chunks.
func (f *File) ChunksCount(_ context.Context) (int, error) {
	return len(f.data), nil
}

// WriteChunk writes data to a chunk.
func (f *File) WriteChunk(_ context.Context, chunkIndex int, start int, end *int, data []byte) (int, error) {
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

// ResizeChunksNb resizes the number of chunks.
func (f *File) ResizeChunksNb(_ context.Context, nb int) error {
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

// ResizeLastChunk resizes the last chunk.
func (f *File) ResizeLastChunk(_ context.Context, size int) (changed int, err error) {
	// Check if the size is valid
	if size < 0 || size > f.chunkSize {
		return 0, storage.ErrInvalidChunkSize
	}

	// Check if there is a chunk to resize
	if len(f.data) == 0 {
		return 0, storage.ErrNoChunk
	}

	lastChunkSize := len(f.data[len(f.data)-1])
	toModify := size - lastChunkSize
	if toModify < 0 {
		// Truncate the last chunk
		f.data[len(f.data)-1] = f.data[len(f.data)-1][:size]
		return toModify, nil
	} else if toModify > 0 {
		// Add data to the last chunk
		f.data[len(f.data)-1] = append(f.data[len(f.data)-1], make([]byte, toModify)...)
		return toModify, nil
	}

	return 0, nil
}

// Size returns the size of the file.
func (f *File) Size(_ context.Context) (int, error) {
	// Check if there is no data
	size := len(f.data)
	if size == 0 {
		return 0, nil
	}

	// Return the count of all chunks except the last one, multiplied by the chunk size
	// + the length of the last chunk
	return (size-1)*f.chunkSize + len(f.data[size-1]), nil
}

// LastChunkSize returns the size of the last chunk.
func (f *File) LastChunkSize(_ context.Context) (int, error) {
	// Check if there is no data
	if len(f.data) == 0 {
		return 0, storage.ErrNoChunk
	}

	return len(f.data[len(f.data)-1]), nil
}
