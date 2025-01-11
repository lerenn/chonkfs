package mem

import (
	"context"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
)

type file struct {
	content []byte
}

func newEmptyFile() *file {
	return &file{
		content: make([]byte, 0),
	}
}

func (f *file) Read(ctx context.Context, off int64) ([]byte, syscall.Errno) {
	return f.content[off:len(f.content)], fs.OK
}

func (f *file) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	// Get everything before off
	content := f.content[:off]

	// Add the data
	content = append(content, data...)

	// Save the data
	f.content = content

	return uint32(len(data)), fs.OK
}
