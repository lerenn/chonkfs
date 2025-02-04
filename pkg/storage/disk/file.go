package disk

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/lerenn/chonkfs/pkg/storage"
)

const (
	defaultFileMode = os.FileMode(0755)
)

var _ storage.File = (*file)(nil)

type fileOptions struct {
	Underlayer    storage.File
	ChunkNb       int
	LastChunkSize int
}

type file struct {
	path       string
	chunkSize  int
	underlayer storage.File
}

func newFile(path string, chunkSize int, opts *fileOptions) (*file, error) {
	f := &file{
		path:      path,
		chunkSize: chunkSize,
	}

	if opts == nil {
		return f, nil
	}

	// Set the underlayer
	f.underlayer = opts.Underlayer

	// Create the chunks if specified
	if opts.ChunkNb > 0 {
		if err := f.addChunks(0, opts.ChunkNb); err != nil {
			return nil, err
		}
	}

	// Change the last chunk size if specified

	return f, nil
}

// Underlayer returns the underlayer file.
func (f *file) Underlayer() storage.File {
	return f.underlayer
}

func (f *file) chunkPath(chunkIndex int) string {
	name := fmt.Sprintf("%d.dat", chunkIndex)
	return path.Join(f.path, name)
}

// Info returns the file info.
func (f *file) Info(_ context.Context) (storage.FileInfo, error) {
	return readFileInfo(f.path)
}

// ReadChunk reads data from a chunk.
func (f *file) ReadChunk(ctx context.Context, chunkIndex int, data []byte, offset int) (int, error) {
	// Check if the parameters are valid
	if err := f.checkReadWriteChunkParams(ctx, chunkIndex, data, offset); err != nil {
		return 0, err
	}

	// Read the data
	of, err := os.Open(f.chunkPath(chunkIndex))
	if err != nil {
		return 0, fmt.Errorf("%w: chunk %d for %q", err, chunkIndex, f.path)
	}
	defer of.Close()

	n, err := of.ReadAt(data, int64(offset))
	if err != nil {
		return n, fmt.Errorf("%w: chunk %d for %q", err, chunkIndex, f.path)
	}

	return n, nil
}

func (f *file) checkReadWriteChunkParams(ctx context.Context, chunkIndex int, data []byte, offset int) error {
	// Get the current number of chunks
	info, err := f.Info(ctx)
	if err != nil {
		return err
	}

	// Check if the chunk index is valid
	if chunkIndex < 0 || chunkIndex >= info.ChunksCount {
		return storage.ErrInvalidChunkNb
	}

	// Check if the start is valid
	if offset < 0 ||
		(chunkIndex != info.ChunksCount-1 && offset >= info.ChunkSize) ||
		(chunkIndex == info.ChunksCount-1 && offset >= info.LastChunkSize) {
		return fmt.Errorf("%w: start is %d", storage.ErrInvalidStartOffset, offset)
	}

	// Check if the end is valid
	end := offset + len(data)
	if end < 0 ||
		(chunkIndex != info.ChunksCount-1 && end >= info.ChunkSize) ||
		(chunkIndex == info.ChunksCount-1 && end >= info.LastChunkSize) {
		return fmt.Errorf("%w: end is %d", storage.ErrInvalidEndOffset, end)
	}

	return nil
}

// WriteChunk writes data to a chunk.
func (f *file) WriteChunk(ctx context.Context, chunkIndex int, data []byte, offset int) (int, error) {
	// Check if the parameters are valid
	if err := f.checkReadWriteChunkParams(ctx, chunkIndex, data, offset); err != nil {
		return 0, err
	}

	// Write it in the underlayer, if there is one
	if u := f.Underlayer(); u != nil {
		if w, err := u.WriteChunk(ctx, chunkIndex, data, offset); err != nil {
			return w, err
		}
	}

	// Write the data
	of, err := os.OpenFile(f.chunkPath(chunkIndex), os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("%w: chunk %d for %q", err, chunkIndex, f.path)
	}
	defer of.Close()

	_, err = of.WriteAt(data, int64(offset))
	if err != nil {
		return 0, fmt.Errorf("%w: chunk %d for %q", err, chunkIndex, f.path)
	}

	return len(data), nil
}

func (f *file) addChunks(oldSize, newSize int) error {
	for i := oldSize; i < newSize; i++ {
		if err := os.WriteFile(f.chunkPath(i), make([]byte, f.chunkSize), defaultFileMode); err != nil {
			return fmt.Errorf("%w: chunk %d for %q", err, i, f.path)
		}
	}
	return nil
}

func (f *file) truncateChunks(oldSize, newSize int) error {
	for i := newSize; i < oldSize; i++ {
		if err := os.Remove(f.chunkPath(i)); err != nil {
			return fmt.Errorf("%w: chunk %d for %q", err, i, f.path)
		}
	}
	return nil
}

// ResizeChunksNb resizes the number of chunks.
func (f *file) ResizeChunksNb(ctx context.Context, nb int) error {
	// Check if the number of chunks is valid
	if nb < 0 {
		return storage.ErrInvalidChunkNb
	}

	// If there is an underlayer, apply the resize here first
	if u := f.Underlayer(); u != nil {
		if err := u.ResizeChunksNb(ctx, nb); err != nil {
			return err
		}
	}

	// Get the current number of chunks
	info, err := f.Info(ctx)
	if err != nil {
		return err
	}

	// Apply the resizing
	if nb > info.ChunksCount {
		// Check if the last chunk is full
		if info.ChunksCount > 0 && info.LastChunkSize != f.chunkSize {
			return storage.ErrLastChunkNotFull
		}

		// Add chunks
		if err := f.addChunks(info.ChunksCount, nb); err != nil {
			return err
		}
	} else if nb < info.ChunksCount {
		// Remove chunks
		if err := f.truncateChunks(info.ChunksCount, nb); err != nil {
			return err
		}
	}

	// Update the file info
	info.ChunksCount = nb
	info.LastChunkSize = f.chunkSize
	return writeFileInfo(info, f.path)
}

// ResizeLastChunk resizes the last chunk.
func (f *file) ResizeLastChunk(ctx context.Context, size int) (changed int, err error) {
	// Check if the size is valid
	if size < 0 || size > f.chunkSize {
		return 0, storage.ErrInvalidChunkSize
	}

	// Get the current number of chunks
	info, err := f.Info(ctx)
	if err != nil {
		return 0, err
	}

	// Check if there is a chunk to resize
	if info.ChunksCount == 0 {
		return 0, storage.ErrNoChunk
	}

	// Apply the resize to the underlayer if there is one
	if u := f.Underlayer(); u != nil {
		if c, err := u.ResizeLastChunk(ctx, size); err != nil {
			return c, err
		}
	}

	toModify := size - info.LastChunkSize
	if toModify < 0 {
		// Truncate the last chunk
		os.Truncate(f.chunkPath(info.ChunksCount-1), int64(size))
	} else if toModify > 0 {
		// Open the last chunk
		of, err := os.OpenFile(f.chunkPath(info.ChunksCount-1), os.O_WRONLY, 0644)
		if err != nil {
			return 0, fmt.Errorf("%w: last chunk for %q", err, f.path)
		}
		defer of.Close()

		// Add data to the last chunk
		_, err = of.WriteAt(make([]byte, toModify), int64(info.LastChunkSize))
		if err != nil {
			return 0, fmt.Errorf("%w: last chunk for %q", err, f.path)
		}
	}

	return toModify, nil
}

// Size returns the size of the file.
func (f *file) Size(_ context.Context) (int, error) {
	return 0, fmt.Errorf("not implemented")
}
