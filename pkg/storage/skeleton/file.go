package skeleton

import (
	"context"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.File = (*file)(nil)

type fileOptions struct {
	Underlayer storage.File
}

type file struct {
	chunkSize  int
	underlayer storage.File
}

func newFile(chunkSize int, opts *fileOptions) *file {
	f := &file{
		chunkSize: chunkSize,
	}

	if opts != nil {
		f.underlayer = opts.Underlayer
	}

	return f
}

// Underlayer returns the underlayer file.
func (f *file) Underlayer() storage.File {
	return f.underlayer
}

// Info returns the file info.
func (f *file) Info(_ context.Context) (storage.FileInfo, error) {
	return storage.FileInfo{}, fmt.Errorf("not implemented")
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

// Size returns the size of the file.
func (f *file) Size(_ context.Context) (int, error) {
	return 0, fmt.Errorf("not implemented")
}
