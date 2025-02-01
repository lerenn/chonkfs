package skeleton

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.File = (*file)(nil)

type file struct{}

func newFile(_ int) *file {
	return &file{}
}

// Underlayer returns the underlayer file.
func (f *file) Underlayer() storage.File {
	return nil
}

// Info returns the file info.
func (f *file) Info(_ context.Context) (storage.FileInfo, error) {
	return storage.FileInfo{}, nil
}

// ReadChunk reads _ from a chunk.
func (f *file) ReadChunk(_ context.Context, _ int, _ []byte, _ int, _ *int) (int, error) {
	return 0, nil
}

// ChunksCount returns the number of chunks.
func (f *file) ChunksCount(_ context.Context) (int, error) {
	return 0, nil
}

// WriteChunk writes _ to a chunk.
func (f *file) WriteChunk(_ context.Context, _ int, _ int, _ *int, _ []byte) (int, error) {
	return 0, nil
}

// ResizeChunksNb resizes the number of chunks.
func (f *file) ResizeChunksNb(_ context.Context, _ int) error {
	return nil
}

// ResizeLastChunk resizes the last chunk.
func (f *file) ResizeLastChunk(_ context.Context, _ int) (changed int, err error) {
	return 0, nil
}

// Size returns the size of the file.
func (f *file) Size(_ context.Context) (int, error) {
	return 0, nil
}

// Last_ returns the size of the last chunk.
func (f *file) LastChunkSize(_ context.Context) (int, error) {
	return 0, nil
}
