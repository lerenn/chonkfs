package layer

import (
	"context"
	"errors"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.File = (*file)(nil)

type file struct {
	chunkSize  int
	backend    storage.File
	underlayer storage.File
}

func newFile(backend storage.File, underlayer storage.File, info info.File) *file {
	return &file{
		chunkSize:  info.ChunkSize,
		backend:    backend,
		underlayer: underlayer,
	}
}

// GetInfo returns the file info.
func (f *file) GetInfo(ctx context.Context) (info.File, error) {
	if fileInfo, err := f.backend.GetInfo(ctx); err == nil {
		return fileInfo, nil
	} else if !errors.Is(err, storage.ErrFileNotFound) {
		return info.File{}, fmt.Errorf("%w: %w", storage.ErrStorage, err)
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
func (f *file) ResizeChunksNb(ctx context.Context, size int) error {
	if err := f.underlayer.ResizeChunksNb(ctx, size); err != nil {
		return err
	}

	return f.backend.ResizeChunksNb(ctx, size)
}

// ResizeLastChunk resizes the last chunk.
func (f *file) ResizeLastChunk(ctx context.Context, size int) (changed int, err error) {
	if _, err := f.underlayer.ResizeLastChunk(ctx, size); err != nil {
		return 0, err
	}

	return f.backend.ResizeLastChunk(ctx, size)
}
