package chonker

import (
	"context"
	"io"
	"log"

	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ File = (*file)(nil)

type fileOption func(fl *file)

// WithFileLogger is an option to set the logger of a file.
//
//nolint:revive
func WithFileLogger(logger *log.Logger) fileOption {
	return func(fl *file) {
		fl.logger = logger
	}
}

type file struct {
	underlayer storage.File
	chunkSize  int

	opts   []fileOption
	logger *log.Logger
}

// NewFile creates a new file.
func NewFile(
	_ context.Context,
	s storage.File,
	chunkSize int,
	opts ...fileOption,
) (File, error) {
	// Create a default file
	f := &file{
		underlayer: s,
		chunkSize:  chunkSize,
		opts:       opts,
		logger:     log.New(io.Discard, "", 0),
	}

	// Apply options
	for _, opt := range opts {
		opt(f)
	}

	return f, nil
}

// GetAttributes returns the attributes of the file.
func (f *file) GetAttributes(ctx context.Context) (FileAttributes, error) {
	size, err := f.Size(ctx)
	if err != nil {
		return FileAttributes{}, err
	}

	return FileAttributes{
		Size: size,
	}, nil
}

// SetAttributes sets the attributes of the file.
func (f *file) SetAttributes(_ context.Context, _ FileAttributes) error {
	// Nothing to do (yet)
	return nil
}

// Read reads the file at the given offset.
func (f *file) Read(ctx context.Context, dest []byte, off int) ([]byte, error) {
	return f.readAccrossChunks(ctx, dest, off)
}

// TODO: Refactor this function
//
//nolint:cyclop
func (f *file) readAccrossChunks(ctx context.Context, dest []byte, off int) ([]byte, error) {
	// Get the total number of chunks
	totalChunk, err := f.underlayer.ChunksCount(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the offset is valid
	chunkNb := off / f.chunkSize
	if chunkNb >= totalChunk {
		return []byte{}, nil
	} else if chunkNb == totalChunk-1 {
		ls, err := f.underlayer.LastChunkSize(ctx)
		if err != nil {
			return nil, err
		}

		if off%f.chunkSize >= ls {
			return []byte{}, nil
		}
	}

	// Loop across chunks
	read := 0
	for ; read < len(dest) && chunkNb < totalChunk; chunkNb++ {
		if read == 0 {
			r, err := f.underlayer.ReadChunk(ctx, chunkNb, dest, off%f.chunkSize, nil)
			if err != nil {
				return nil, err
			}
			read += r
		} else {
			r, err := f.underlayer.ReadChunk(ctx, chunkNb, dest[read:], 0, nil)
			if err != nil {
				return nil, err
			}
			read += r
		}
	}

	// Check if the read is less than the destination size
	if read < len(dest) {
		return dest[:read], nil
	}

	return dest, nil
}

// Write writes the data at the given offset.
func (f *file) Write(ctx context.Context, data []byte, off int, opts WriteOptions) (written int, err error) {
	// Check if there is enough space, and allocate what's missing
	if err := f.resizeChunks(ctx, off+len(data)); err != nil {
		return 0, err
	}

	// Check if we need to append
	if opts.Append {
		if err := f.append(ctx, data); err != nil {
			return 0, err
		}

		return len(data), nil
	}

	// Write the data
	written, err = f.writeAccrossChunks(ctx, data, off)
	if err != nil {
		return 0, err
	}

	// Check if truncate is needed
	if opts.Truncate {
		if err := f.Truncate(ctx, off+len(data)); err != nil {
			return 0, err
		}
	}

	return written, nil
}

func (f *file) append(ctx context.Context, data []byte) error {
	oldSize, err := f.Size(ctx)
	if err != nil {
		return err
	}

	// Add missing chunks
	if err := f.resizeChunks(ctx, oldSize+len(data)); err != nil {
		return err
	}

	// Write the data at the end
	_, err = f.writeAccrossChunks(ctx, data, oldSize)
	return err
}

// Truncate truncates the file to the given size.
func (f *file) Truncate(ctx context.Context, newSize int) error {
	// Check if we need to truncate
	oldSize, err := f.Size(ctx)
	if err != nil {
		return err
	} else if newSize >= oldSize {
		return nil
	}

	// Check if we need to truncate the last chunk, and remove the chunks after the new last one
	// NOTE: -1 is used to avoid the case where the last chunk is full
	if (oldSize-1)/f.chunkSize != (newSize-1)/f.chunkSize {
		lastChunkNb := newSize / f.chunkSize
		if err := f.underlayer.ResizeChunksNb(ctx, lastChunkNb); err != nil {
			return err
		}
	}

	// Truncate the last chunk if needed
	partialLastChunkSize := newSize % f.chunkSize
	if partialLastChunkSize > 0 {
		if _, err := f.underlayer.ResizeLastChunk(ctx, partialLastChunkSize); err != nil {
			return err
		}
	}

	return nil
}

// Size returns the size of the file.
func (f *file) Size(ctx context.Context) (int, error) {
	return f.underlayer.Size(ctx)
}

func (f *file) resizeChunks(ctx context.Context, newSize int) error {
	// Check if there is enough space, and allocate what's missing
	oldSize, err := f.Size(ctx)
	if err != nil {
		return err
	} else if newSize <= oldSize {
		return err
	}

	// Check if the last chunk is not full
	if oldSize%f.chunkSize != 0 {
		// Make it full or less
		if err := f.makeLastChunkFullOrLess(ctx, oldSize, &newSize); err != nil {
			return err
		}

		// If there is nothing else to do, return
		if newSize == oldSize {
			return nil
		}
	}

	// Get the total number of full chunks and the size of the last chunk
	fullChunksToAdd := newSize / f.chunkSize
	partialLastChunkToAdd := newSize % f.chunkSize
	if partialLastChunkToAdd > 0 {
		fullChunksToAdd++
	}

	// Add the missing chunks
	if err := f.underlayer.ResizeChunksNb(ctx, fullChunksToAdd); err != nil {
		return err
	}

	// Truncate the last chunk if needed
	if partialLastChunkToAdd > 0 {
		if _, err := f.underlayer.ResizeLastChunk(ctx, partialLastChunkToAdd); err != nil {
			return err
		}
	}

	return nil
}

func (f *file) makeLastChunkFullOrLess(ctx context.Context, oldSize int, newSize *int) error {
	// Get chunk count
	chunkNb, err := f.underlayer.ChunksCount(ctx)
	if err != nil {
		return err
	}

	// Ask for a full resize
	resize := f.chunkSize

	// However, if the total size is smaller than the chunk size, resize to the total size
	if *newSize-oldSize < f.chunkSize {
		resize = *newSize - (chunkNb-1)*f.chunkSize
	}

	// Resize the last chunk
	added, err := f.underlayer.ResizeLastChunk(ctx, resize)
	if err != nil {
		return err
	}
	*newSize -= added

	return nil
}

func (f *file) writeAccrossChunks(ctx context.Context, data []byte, off int) (written int, err error) {
	size, err := f.Size(ctx)
	if err != nil {
		return 0, err
	}

	for chunkNb := off / f.chunkSize; written < len(data) && chunkNb < size; chunkNb++ {
		if written == 0 {
			w, err := f.underlayer.WriteChunk(ctx, chunkNb, off%f.chunkSize, nil, data)
			if err != nil {
				return 0, err
			}
			written += w
		} else {
			w, err := f.underlayer.WriteChunk(ctx, chunkNb, 0, nil, data[written:])
			if err != nil {
				return 0, err
			}
			written += w
		}
	}

	return written, nil
}

// Sync saves the file to the storage.
func (f *file) Sync(_ context.Context) error {
	// TODO: Save to a embedded backend
	return nil
}
