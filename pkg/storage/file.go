package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage/backend"
)

var _ File = (*file)(nil)

type fileOptions struct {
	Underlayer File
}

type file struct {
	chunkSize  int
	backend    backend.File
	underlayer File
}

func newFile(backend backend.File, chunkSize int, opts *fileOptions) *file {
	f := &file{
		chunkSize: chunkSize,
		backend:   backend,
	}

	if opts != nil {
		f.underlayer = opts.Underlayer
	}

	return f
}

// Underlayer returns the underlayer file.
func (f *file) Underlayer() File {
	return f.underlayer
}

// ChunkSize returns the chunk size.
func (f *file) ChunkSize() int {
	return f.chunkSize
}

// GetInfo returns the file info.
func (f *file) GetInfo(ctx context.Context) (info.File, error) {
	if fileInfo, err := f.backend.GetInfo(ctx); err == nil {
		return fileInfo, nil
	} else if !errors.Is(err, backend.ErrNotFound) {
		return info.File{}, fmt.Errorf("%w: %w", ErrStorage, err)
	}

	if f.underlayer == nil {
		return info.File{}, ErrFileNotExists
	}

	return f.underlayer.GetInfo(ctx)
}

// ReadChunk reads _ from a chunk.
func (f *file) ReadChunk(_ context.Context, _ int, _ []byte, _ int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

// WriteChunk writes _ to a chunk.
func (f *file) WriteChunk(_ context.Context, _ int, _ []byte, _ int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

// ResizeChunksNb resizes the number of chunks.
func (f *file) ResizeChunksNb(_ context.Context, _ int) error {
	return fmt.Errorf("not implemented")
}

// ResizeLastChunk resizes the last chunk.
func (f *file) ResizeLastChunk(_ context.Context, _ int) (changed int, err error) {
	return 0, fmt.Errorf("not implemented")
}
