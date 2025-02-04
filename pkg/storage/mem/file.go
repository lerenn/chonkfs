package mem

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.File = (*file)(nil)

type fileOptions struct {
	Underlayer    storage.File
	ChunkNb       int
	LastChunkSize int
}

type file struct {
	data       [][]byte
	chunkSize  int
	underlayer storage.File
}

func newFile(chunkSize int, opts *fileOptions) *file {
	f := &file{
		data:      make([][]byte, 0),
		chunkSize: chunkSize,
	}

	if opts == nil {
		return f
	}

	// Set the underlayer
	f.underlayer = opts.Underlayer

	// Create the chunks if specified
	if opts.ChunkNb > 0 {
		for i := 0; i < opts.ChunkNb; i++ {
			f.data = append(f.data, make([]byte, chunkSize))
		}
	}

	// Change the last chunk size if specified
	if opts.ChunkNb > 0 && opts.LastChunkSize > 0 {
		f.data[len(f.data)-1] = make([]byte, opts.LastChunkSize)
	}

	return f
}

// Underlayer returns the underlayer file.
func (f *file) Underlayer() storage.File {
	return f.underlayer
}

// Info returns the file info.
func (f *file) Info(_ context.Context) (storage.FileInfo, error) {
	// Check if there is no data
	lastChunkSize := 0
	if len(f.data) > 0 {
		lastChunkSize = len(f.data[len(f.data)-1])
	}

	return storage.FileInfo{
		ChunkSize:     f.chunkSize,
		ChunksCount:   len(f.data),
		LastChunkSize: lastChunkSize,
	}, nil
}

// ReadChunk reads data from a chunk.
func (f *file) ReadChunk(_ context.Context, chunkIndex int, data []byte, offset int) (int, error) {
	// Check if the chunk index is valid
	if chunkIndex < 0 || chunkIndex >= len(f.data) {
		return 0, storage.ErrInvalidChunkNb
	}

	// Check if the start is valid
	if offset < 0 || offset >= len(f.data[chunkIndex]) {
		return 0, fmt.Errorf("%w: start is %d", storage.ErrInvalidStartOffset, offset)
	}

	// Check if the end is valid
	end := offset + len(data)
	if end < 0 || end > len(f.data[chunkIndex]) {
		return 0, fmt.Errorf("%w: end is %d", storage.ErrInvalidEndOffset, end)
	}

	// Read the data
	return copy(data, f.data[chunkIndex][offset:end]), nil
}

func (f *file) checkWriteChunkParams(chunkIndex int, data []byte, offset int) error {
	// Check if the chunk index is valid
	if chunkIndex < 0 || chunkIndex >= len(f.data) {
		return storage.ErrInvalidChunkNb
	}

	// Check if the start is valid
	if offset < 0 || offset >= len(f.data[chunkIndex]) {
		return fmt.Errorf("%w: offset is %d", storage.ErrInvalidStartOffset, offset)
	}

	// Check if the end is valid
	end := offset + len(data)
	if end < 0 || end > len(f.data[chunkIndex]) {
		return fmt.Errorf("%w: end is %d", storage.ErrInvalidEndOffset, end)
	}

	return nil
}

// WriteChunk writes data to a chunk.
func (f *file) WriteChunk(ctx context.Context, chunkIndex int, data []byte, offset int) (int, error) {
	// Check params
	if err := f.checkWriteChunkParams(chunkIndex, data, offset); err != nil {
		return 0, err
	}

	// Write it in the underlayer, if there is one
	if u := f.Underlayer(); u != nil {
		if w, err := u.WriteChunk(ctx, chunkIndex, data, offset); err != nil {
			return w, err
		}
	}

	// Write the data
	return copy(f.data[chunkIndex][offset:offset+len(data)], data), nil
}

// ResizeChunksNb resizes the number of chunks.
func (f *file) ResizeChunksNb(ctx context.Context, nb int) error {
	// Check if the number of chunks is valid
	if nb < 0 {
		return storage.ErrInvalidChunkNb
	}

	// If there is an underlayer
	if u := f.Underlayer(); u != nil {
		// Resize the underlayer
		if err := u.ResizeChunksNb(ctx, nb); err != nil {
			return err
		}

		// Pull the data from the underlayer
		for i := len(f.data); i < nb; i++ {
			data := make([]byte, f.chunkSize)
			u.ReadChunk(ctx, i, data, 0)
			f.data = append(f.data, data)
		}

		return nil
	}

	// If there is no underlayer, apply the resizing directly
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
func (f *file) ResizeLastChunk(ctx context.Context, size int) (changed int, err error) {
	// Check if the size is valid
	if size < 0 || size > f.chunkSize {
		return 0, storage.ErrInvalidChunkSize
	}

	// Check if there is a chunk to resize
	if len(f.data) == 0 {
		return 0, storage.ErrNoChunk
	}

	// Apply the resize to the underlayer if there is one
	if u := f.Underlayer(); u != nil {
		if c, err := u.ResizeLastChunk(ctx, size); err != nil {
			return c, err
		}
	}

	lastChunkSize := len(f.data[len(f.data)-1])
	toModify := size - lastChunkSize
	if toModify < 0 {
		// Truncate the last chunk
		f.data[len(f.data)-1] = f.data[len(f.data)-1][:size]
	} else if toModify > 0 {
		// Add data to the last chunk
		f.data[len(f.data)-1] = append(f.data[len(f.data)-1], make([]byte, toModify)...)
	}

	return toModify, nil
}

// Size returns the size of the file.
func (f *file) Size(_ context.Context) (int, error) {
	// Check if there is no data
	size := len(f.data)
	if size == 0 {
		return 0, nil
	}

	// Return the count of all chunks except the last one, multiplied by the chunk size
	// + the length of the last chunk
	return (size-1)*f.chunkSize + len(f.data[size-1]), nil
}
