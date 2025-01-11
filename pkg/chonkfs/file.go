package chonkfs

import (
	"context"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/backend"
)

// Capabilities that the file struct should implements
var (
	_ fs.FileReader = (*file)(nil)
	_ fs.FileWriter = (*file)(nil)

	_ fs.InodeEmbedder = (*file)(nil)

	_ fs.NodeOpener    = (*file)(nil)
	_ fs.NodeSetattrer = (*file)(nil)
)

type file struct {
	backendFile backend.File

	fs.Inode

	// implementers.NodeImplementer
	// implementers.FileImplementer
}

func (fl *file) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	debugf("file.Read\n")

	// Get content from file
	content, errno := fl.backendFile.Read(ctx, off)
	if errno != fs.OK {
		return nil, errno
	}

	return fuse.ReadResultData(content), fs.OK
}

func (fl *file) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	debugf("file.Open\n")

	// Nothing to do for the moment.
	// Please open a ticket if needed.

	return fl, fuse.FOPEN_NOFLUSH | fuse.FOPEN_DIRECT_IO, fs.OK
}

func (fl *file) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	debugf("file.Setattr\n")

	// Nothing to do for the moment.
	// Please open a ticket if needed.

	return fs.OK
}

func (fl *file) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	debugf("file.Write\n")

	// Write content to file
	return fl.backendFile.Write(ctx, data, off)
}
