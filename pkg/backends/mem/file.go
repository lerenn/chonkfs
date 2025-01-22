package mem

import (
	"context"
	"io"
	"log"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/backends"
)

var _ backends.File = (*file)(nil)

type FileOption func(fl *file)

func WithFileLogger(logger *log.Logger) FileOption {
	return func(fl *file) {
		fl.logger = logger
	}
}

type file struct {
	data      [][]byte
	chunkSize int
	logger    *log.Logger
}

func newFile(chunkSize int, opts ...FileOption) *file {
	// Create a default file
	f := &file{
		data:      make([][]byte, 0),
		chunkSize: chunkSize,
		logger:    log.New(io.Discard, "", 0),
	}

	// Apply options
	for _, opt := range opts {
		opt(f)
	}

	return f
}

func (f *file) GetAttributes(ctx context.Context) (fuse.Attr, error) {
	size, _ := f.Size(ctx)

	return fuse.Attr{
		Size: uint64(size),
	}, nil
}

func (f *file) SetAttributes(ctx context.Context, in *fuse.SetAttrIn) error {
	// TODO
	return nil
}

func (f *file) Read(ctx context.Context, data []byte, off int) error {
	f.readAccrossChunks(data, off)
	return nil
}

func (f *file) readAccrossChunks(data []byte, off int) {
	for chunkNb, read := off/f.chunkSize, 0; read < len(data) && chunkNb < len(f.data); chunkNb++ {
		if read == 0 {
			read += copy(data, f.data[chunkNb][off%f.chunkSize:])
		} else {
			read += copy(data[read:], f.data[chunkNb])
		}
	}
}

func (f *file) Write(ctx context.Context, data []byte, off int, opts backends.WriteOptions) (written int, errno error) {
	// Check if there is enough space, and allocate what's missing
	f.addMissingChunks(ctx, off+len(data))

	// Check if we need to append
	if opts.Append {
		f.append(ctx, data)
	} else {
		// Write the data
		f.writeAccrossChunks(data, off)

		// Check if truncate is needed
		if opts.Truncate {
			if err := f.Truncate(ctx, off+len(data)); err != nil {
				return 0, err
			}
		}
	}

	return len(data), nil
}

func (f *file) append(ctx context.Context, data []byte) {
	oldSize, _ := f.Size(ctx)

	// Add missing chunks
	f.addMissingChunks(ctx, oldSize+len(data))

	// Write the data at the end
	f.writeAccrossChunks(data, oldSize)
}

func (f *file) Truncate(ctx context.Context, newSize int) error {
	// Check if we need to truncate
	oldSize, _ := f.Size(ctx)
	if newSize >= oldSize {
		return nil
	}

	// Check if we need to truncate the last chunk
	// NOTE: -1 is used to avoid the case where the last chunk is full
	if (oldSize-1)/f.chunkSize == (newSize-1)/f.chunkSize {
		newSize -= f.truncateLastChunk(newSize % f.chunkSize)
		if newSize == 0 {
			return nil
		}
	}

	// Compute the new last chunk number
	lastChunkNb := newSize / f.chunkSize
	partialLastChunkSize := newSize % f.chunkSize

	// Remove all the chunks after the new last one
	f.data = f.data[:lastChunkNb]

	// Truncate the last chunk if needed
	if partialLastChunkSize > 0 {
		f.data[lastChunkNb] = f.data[lastChunkNb][:partialLastChunkSize]
	}

	return nil
}

func (f *file) truncateLastChunk(newSize int) (truncated int) {
	oldSize := len(f.data[len(f.data)-1])
	if oldSize > newSize {
		// The last chunk size is bigger than the new size: truncate it then stop
		f.data[len(f.data)-1] = f.data[len(f.data)-1][:newSize]
		return newSize
	} else if oldSize == newSize {
		// The last chunk size is equal to the new size: remove it then stop
		f.data = f.data[:len(f.data)-1]
		return newSize
	} else {
		// The last chunk size is smaller than the new size: remove it and continue
		f.data = f.data[:len(f.data)-1]
		return newSize - oldSize
	}
}

func (f *file) Size(ctx context.Context) (int, error) {
	if len(f.data) == 0 {
		return 0, nil
	}

	completeChunksNb := len(f.data) - 1
	lastChunk := f.data[len(f.data)-1]
	return completeChunksNb*f.chunkSize + len(lastChunk), nil
}

func (f *file) addMissingChunks(ctx context.Context, total int) {
	// Check if there is enough space, and allocate what's missing
	size, _ := f.Size(ctx)
	if int(total) <= size {
		return
	}

	// Check if the last chunk is not full
	lastChunkNb := len(f.data) - 1
	if len(f.data) > 0 && len(f.data[lastChunkNb]) < f.chunkSize {
		// Add the missing space to make it full and substract it from the total
		f.data[lastChunkNb] = append(f.data[lastChunkNb], make([]byte, f.chunkSize-len(f.data[lastChunkNb]))...)
		total -= f.chunkSize - len(f.data[lastChunkNb])
	}

	// Get the total number of full chunks and the size of the last chunk
	fullChunksToAdd := total / f.chunkSize
	partialLastChunkToAdd := total % f.chunkSize

	// Add the missing chunks
	for i := 0; i < fullChunksToAdd; i++ {
		f.data = append(f.data, make([]byte, f.chunkSize))
	}

	// Add the last chunk if needed
	if partialLastChunkToAdd > 0 {
		f.data = append(f.data, make([]byte, partialLastChunkToAdd))
	}
}

func (f *file) writeAccrossChunks(data []byte, off int) {
	for chunkNb, written := off/f.chunkSize, 0; written < len(data) && chunkNb < len(f.data); chunkNb++ {
		if written == 0 {
			written += copy(f.data[chunkNb][off%f.chunkSize:], data[written:])
		} else {
			written += copy(f.data[chunkNb], data[written:])
		}
	}
}

func (f *file) Sync(ctx context.Context) error {
	// TODO: Save to a embedded backend
	return nil
}
