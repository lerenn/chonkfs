package mem

import (
	"context"
	"slices"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/backends"
)

var _ backends.File = (*file)(nil)

type file struct {
	attr    fuse.Attr
	cache   []byte
	content []byte
}

func newEmptyFile() *file {
	return &file{
		content: make([]byte, 0),
		cache:   make([]byte, 0),
	}
}

func (f *file) GetAttributes(ctx context.Context) (fuse.Attr, error) {
	return f.attr, nil
}

func (f *file) SetAttributes(ctx context.Context, in *fuse.SetAttrIn) error {
	// TODO
	return nil
}

func (f *file) Read(ctx context.Context, start, end int) ([]byte, error) {
	// Check that the end is after the start
	if end < start {
		return nil, backends.ErrReadEndBeforeReadStart
	}

	// If there is no cache, populate it with the content
	if len(f.cache) == 0 {
		f.cache = slices.Clone(f.content)
	}

	// Check that the offset is within the cache
	if end >= len(f.cache) {
		end = len(f.cache) - 1
	}

	return f.cache[start:end], nil
}

func (f *file) WriteCache(ctx context.Context, data []byte, off int) (written int, err error) {
	// Check if there is enough space, and allocate what's missing
	if int(off) > len(f.cache) {
		f.cache = append(f.cache, make([]byte, int(off)-(len(f.cache)-1))...)
	}

	// Get everything before off and after the data's end
	before, after := f.cache[:off], f.cache[off+len(data):]

	// Concatenate everything
	cache := append(before, data...)
	cache = append(cache, after...)

	// Save the data
	f.cache = cache

	return len(data), nil
}

func (f *file) Sync(ctx context.Context) error {
	// Write to content
	f.content = slices.Clone(f.cache)

	// Update attributes based on content
	f.attr.Size = uint64(len(f.content))

	return nil
}
