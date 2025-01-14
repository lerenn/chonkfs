package mem

import (
	"context"
	"slices"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type file struct {
	attr    fuse.Attr
	cache   []byte
	content []byte
}

func newEmptyFile() *file {
	return &file{
		content: make([]byte, 0),
	}
}

func (f *file) GetAttributes(ctx context.Context, attr *fuse.Attr) (errno syscall.Errno) {
	*attr = f.attr
	return fs.OK
}

func (f *file) SetAttributes(ctx context.Context, in *fuse.SetAttrIn) syscall.Errno {
	// TODO
	return fs.OK
}

func (f *file) Read(ctx context.Context, off int64) ([]byte, syscall.Errno) {
	if len(f.cache) == 0 {
		f.cache = slices.Clone(f.content)
	}

	return f.content[off:len(f.content)], fs.OK
}

func (f *file) WriteCache(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	// Get everything before off
	cache := f.cache[:off]

	// Add the data
	cache = append(cache, data...)

	// Save the data
	f.cache = cache

	return uint32(len(data)), fs.OK
}

func (f *file) Sync(ctx context.Context) syscall.Errno {
	// Write to content
	f.content = slices.Clone(f.cache)

	// Update attributes based on content
	f.attr.Size = uint64(len(f.content))

	return fs.OK
}
