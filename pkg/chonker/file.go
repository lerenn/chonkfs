package chonker

import (
	"context"
	"io"
	"log"

	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ File = (*file)(nil)

type FileOption func(fl *file)

func WithFileLogger(logger *log.Logger) FileOption {
	return func(fl *file) {
		fl.logger = logger
	}
}

type file struct {
	storageFile storage.File
	path        string
	chunkSize   int

	opts   []FileOption
	logger *log.Logger
}

func NewFile(
	ctx context.Context,
	s storage.File,
	chunkSize int,
	opts ...FileOption,
) (File, error) {
	// Create a default file
	f := &file{
		storageFile: s,
		chunkSize:   chunkSize,
		opts:        opts,
		logger:      log.New(io.Discard, "", 0),
	}

	// Apply options
	for _, opt := range opts {
		opt(f)
	}

	return f, nil
}

func (f *file) GetAttributes(ctx context.Context) (FileAttributes, error) {
	size, err := f.Size(ctx)
	if err != nil {
		return FileAttributes{}, err
	}

	return FileAttributes{
		Size: size,
	}, nil
}

func (f *file) SetAttributes(ctx context.Context, attr FileAttributes) error {
	// Nothing to do (yet)
	return nil
}

func (f *file) Read(ctx context.Context, dest []byte, off int) ([]byte, error) {
	return f.readAccrossChunks(ctx, dest, off)
}

func (f *file) readAccrossChunks(ctx context.Context, dest []byte, off int) ([]byte, error) {
	totalChunk, err := f.storageFile.ChunksCount(ctx)
	if err != nil {
		return nil, err
	}

	read := 0
	for chunkNb := off / f.chunkSize; read < len(dest) && chunkNb < totalChunk; chunkNb++ {
		if read == 0 {
			r, err := f.storageFile.ReadChunk(ctx, chunkNb, dest, off%f.chunkSize, nil)
			if err != nil {
				return nil, err
			}
			read += r
		} else {
			r, err := f.storageFile.ReadChunk(ctx, chunkNb, dest[read:], 0, nil)
			if err != nil {
				return nil, err
			}
			read += r
		}
	}

	// Reduce dest to the actual read size if needed
	if len(dest) > read {
		dest = dest[:off+read]
	}

	return dest, nil
}

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

		written = len(data)
	} else {
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
		if err := f.storageFile.ResizeChunksNb(ctx, lastChunkNb); err != nil {
			return err
		}
	}

	// Truncate the last chunk if needed
	partialLastChunkSize := newSize % f.chunkSize
	if partialLastChunkSize > 0 {
		if _, err := f.storageFile.ResizeLastChunk(ctx, partialLastChunkSize); err != nil {
			return err
		}
	}

	return nil
}

func (f *file) Size(ctx context.Context) (int, error) {
	return f.storageFile.Size(ctx)
}

func (f *file) resizeChunks(ctx context.Context, total int) error {
	// Check if there is enough space, and allocate what's missing
	size, err := f.Size(ctx)
	if err != nil {
		return err
	} else if int(total) <= size {
		return err
	}

	chunkNb, err := f.storageFile.ChunksCount(ctx)
	if err != nil {
		return err
	}

	// Check if the last chunk is not full, make it full if needed
	if size%f.chunkSize != 0 {
		// Ask for a full resize
		resize := f.chunkSize

		// However, if the total size is smaller than the chunk size, resize to the total size
		if total-size < f.chunkSize {
			resize = total - (chunkNb-1)*f.chunkSize
		}

		// Resize the last chunk
		added, err := f.storageFile.ResizeLastChunk(ctx, resize)
		if err != nil {
			return err
		}
		total -= added

		// If there is nothing else to do, return
		if total == size {
			return nil
		}
	}

	// Get the total number of full chunks and the size of the last chunk
	fullChunksToAdd := total / f.chunkSize
	partialLastChunkToAdd := total % f.chunkSize
	if partialLastChunkToAdd > 0 {
		fullChunksToAdd++
	}

	// Add the missing chunks
	if err := f.storageFile.ResizeChunksNb(ctx, fullChunksToAdd); err != nil {
		return err
	}

	// Truncate the last chunk if needed
	if partialLastChunkToAdd > 0 {
		if _, err := f.storageFile.ResizeLastChunk(ctx, partialLastChunkToAdd); err != nil {
			return err
		}
	}

	return nil
}

func (f *file) writeAccrossChunks(ctx context.Context, data []byte, off int) (written int, err error) {
	size, err := f.Size(ctx)
	if err != nil {
		return 0, err
	}

	for chunkNb := off / f.chunkSize; written < len(data) && chunkNb < size; chunkNb++ {
		if written == 0 {
			w, err := f.storageFile.WriteChunk(ctx, chunkNb, off%f.chunkSize, nil, data)
			if err != nil {
				return 0, err
			}
			written += w
		} else {
			w, err := f.storageFile.WriteChunk(ctx, chunkNb, 0, nil, data[written:])
			if err != nil {
				return 0, err
			}
			written += w
		}
	}

	return written, nil
}

func (f *file) Sync(ctx context.Context) error {
	// TODO: Save to a embedded backend
	return nil
}
