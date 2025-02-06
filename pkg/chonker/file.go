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
	storage   storage.File
	chunkSize int

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
		storage:   s,
		chunkSize: chunkSize,
		opts:      opts,
		logger:    log.New(io.Discard, "", 0),
	}

	// Apply options
	for _, opt := range opts {
		opt(f)
	}

	return f, nil
}

// GetAttributes returns the attributes of the file.
func (f *file) GetAttributes(ctx context.Context) (FileAttributes, error) {
	info, err := f.storage.GetInfo(ctx)
	if err != nil {
		return FileAttributes{}, err
	}

	return FileAttributes{
		Size: info.Size,
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
	// Get info from the underlayer
	info, err := f.storage.GetInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the offset is valid
	chunkNb := off / f.chunkSize
	if chunkNb >= info.ChunksCount {
		return []byte{}, nil
	} else if chunkNb == info.ChunksCount-1 {
		if off%f.chunkSize >= info.LastChunkSize {
			return []byte{}, nil
		}
	}

	// Loop across chunks
	read := 0
	for ; read < len(dest) && chunkNb < info.ChunksCount; chunkNb++ {
		if read == 0 {
			r, err := f.storage.ReadChunk(ctx, chunkNb, dest, off%f.chunkSize)
			if err != nil {
				return nil, err
			}
			read += r
		} else {
			r, err := f.storage.ReadChunk(ctx, chunkNb, dest[read:], 0)
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
	info, err := f.storage.GetInfo(ctx)
	if err != nil {
		return err
	}
	oldSize := info.Size

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
	info, err := f.storage.GetInfo(ctx)
	if err != nil {
		return err
	}

	oldSize := info.Size
	if newSize >= oldSize {
		return nil
	}

	// Check if we need to truncate the last chunk, and remove the chunks after the new last one
	// NOTE: -1 is used to avoid the case where the last chunk is full
	if (oldSize-1)/f.chunkSize != (newSize-1)/f.chunkSize {
		lastChunkNb := newSize / f.chunkSize
		if err := f.storage.ResizeChunksNb(ctx, lastChunkNb); err != nil {
			return err
		}
	}

	// Truncate the last chunk if needed
	partialLastChunkSize := newSize % f.chunkSize
	if partialLastChunkSize > 0 {
		if _, err := f.storage.ResizeLastChunk(ctx, partialLastChunkSize); err != nil {
			return err
		}
	}

	return nil
}

func (f *file) resizeChunks(ctx context.Context, newSize int) error {
	// Check if there is enough space, and allocate what's missing
	info, err := f.storage.GetInfo(ctx)
	if err != nil {
		return err
	}

	oldSize := info.Size
	if newSize <= oldSize {
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
	if err := f.storage.ResizeChunksNb(ctx, fullChunksToAdd); err != nil {
		return err
	}

	// Truncate the last chunk if needed
	if partialLastChunkToAdd > 0 {
		if _, err := f.storage.ResizeLastChunk(ctx, partialLastChunkToAdd); err != nil {
			return err
		}
	}

	return nil
}

func (f *file) makeLastChunkFullOrLess(ctx context.Context, oldSize int, newSize *int) error {
	// Get info
	info, err := f.storage.GetInfo(ctx)
	if err != nil {
		return err
	}

	// Ask for a full resize
	resize := f.chunkSize

	// However, if the total size is smaller than the chunk size, resize to the total size
	if *newSize-oldSize < f.chunkSize {
		resize = *newSize - (info.ChunksCount-1)*f.chunkSize
	}

	// Resize the last chunk
	added, err := f.storage.ResizeLastChunk(ctx, resize)
	if err != nil {
		return err
	}
	*newSize -= added

	return nil
}

func (f *file) writeAccrossChunks(ctx context.Context, data []byte, off int) (written int, err error) {
	info, err := f.storage.GetInfo(ctx)
	if err != nil {
		return 0, err
	}
	size := info.Size

	for chunkNb := off / f.chunkSize; written < len(data) && chunkNb < size; chunkNb++ {
		if written == 0 {
			w, err := f.storage.WriteChunk(ctx, chunkNb, data, off%f.chunkSize)
			if err != nil {
				return 0, err
			}
			written += w
		} else {
			w, err := f.storage.WriteChunk(ctx, chunkNb, data[written:], 0)
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
	// TODO: Save to a embedded backend if the option for direct io is not set
	return nil
}
