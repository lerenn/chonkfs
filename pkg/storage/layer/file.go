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
	upperlayer storage.File
	underlayer storage.File
}

func newFile(upperlayer storage.File, underlayer storage.File, info info.File) *file {
	return &file{
		chunkSize:  info.ChunkSize,
		upperlayer: upperlayer,
		underlayer: underlayer,
	}
}

// GetInfo returns the file info.
func (f *file) GetInfo(ctx context.Context) (info.File, error) {
	if fileInfo, err := f.upperlayer.GetInfo(ctx); err == nil {
		return fileInfo, nil
	} else if !errors.Is(err, storage.ErrFileNotFound) {
		return info.File{}, fmt.Errorf("%w: %w", storage.ErrStorage, err)
	}

	return f.underlayer.GetInfo(ctx)
}

// ReadChunk reads _ from a chunk.
func (f *file) ReadChunk(ctx context.Context, index int, data []byte, offset int) (int, error) {
	read, err := f.upperlayer.ReadChunk(ctx, index, data, offset)
	if err == nil {
		return read, nil
	} else if !errors.Is(err, storage.ErrChunkNotFound) {
		return 0, fmt.Errorf("%w: %w", storage.ErrStorage, err)
	}

	return f.underlayer.ReadChunk(ctx, index, data, offset)
}

// WriteChunk writes _ to a chunk.
func (f *file) WriteChunk(ctx context.Context, index int, data []byte, offset int) (int, error) {
	// Write to underlayer
	rd, err := f.underlayer.WriteChunk(ctx, index, data, offset)
	if err != nil {
		return rd, fmt.Errorf("%w: %w", storage.ErrStorage, err)
	}

	// Try to write upperlayer
	_, err = f.upperlayer.WriteChunk(ctx, index, data, offset)
	if err != nil && !errors.Is(err, storage.ErrChunkNotFound) {
		return rd, fmt.Errorf("%w: %w", storage.ErrStorage, err)
	}

	return rd, nil
}

// ResizeChunksNb resizes the number of chunks.
func (f *file) ResizeChunksNb(ctx context.Context, size int) error {
	if err := f.underlayer.ResizeChunksNb(ctx, size); err != nil {
		return err
	}

	return f.upperlayer.ResizeChunksNb(ctx, size)
}

// ResizeLastChunk resizes the last chunk.
func (f *file) ResizeLastChunk(ctx context.Context, size int) (int, error) {
	// Modify it on underlayer
	changed, err := f.underlayer.ResizeLastChunk(ctx, size)
	if err != nil {
		return 0, err
	}

	// Modify it on upperlayer
	if _, err := f.upperlayer.ResizeLastChunk(ctx, size); err == nil {
		return changed, nil
	} else if !errors.Is(err, storage.ErrChunkNotFound) {
		return 0, fmt.Errorf("%w: %w", storage.ErrStorage, err)
	}

	// Get info
	info, err := f.GetInfo(ctx)
	if err != nil {
		return 0, err
	}

	// Import it from underlayer
	return changed, f.importChunkFromUnderlayer(ctx, info.ChunksCount-1)
}

func (f *file) importChunkFromUnderlayer(ctx context.Context, index int) error {
	// Get the chunk from underlayer
	data := make([]byte, f.chunkSize)
	read, err := f.underlayer.ReadChunk(ctx, index, data, 0)
	if err != nil {
		return err
	}

	// Write the chunk to the upperlayer
	return f.upperlayer.ImportChunk(ctx, index, data[:read])
}

func (f *file) ImportChunk(ctx context.Context, index int, data []byte) error {
	if err := f.underlayer.ImportChunk(ctx, index, data); err != nil {
		return err
	}

	return f.upperlayer.ImportChunk(ctx, index, data)
}
