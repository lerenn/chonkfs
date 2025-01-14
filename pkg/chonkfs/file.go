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
	_ fs.FileFlusher   = (*file)(nil)
	_ fs.FileGetattrer = (*file)(nil)
	_ fs.FileReader    = (*file)(nil)
	_ fs.FileWriter    = (*file)(nil)
	_ fs.FileFsyncer   = (*file)(nil)

	_ fs.InodeEmbedder = (*file)(nil)

	_ fs.NodeOpener    = (*file)(nil)
	_ fs.NodeSetattrer = (*file)(nil)
)

type file struct {
	backendFile backend.File
	name        string

	fs.Inode

	// implementers.NodeImplementer
	// implementers.FileImplementer
}

func (fl *file) Getattr(ctx context.Context, out *fuse.AttrOut) (errno syscall.Errno) {
	debugf("file[name=%q].Getattr\n", fl.name)
	return fl.backendFile.GetAttributes(ctx, &out.Attr)
}

func (fl *file) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	debugf("file[name=%q].Read\n", fl.name)

	// Get content from file
	content, errno := fl.backendFile.Read(ctx, off)
	if errno != fs.OK {
		return nil, errno
	}

	return fuse.ReadResultData(content), fs.OK
}

func (fl *file) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	debugf("file[name=%q].Open\n", fl.name)

	// Nothing to do for the moment.
	// Please open a ticket if needed.

	return fl, fuse.FOPEN_DIRECT_IO, fs.OK
}

func (fl *file) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	debugf("file[name=%q].Write\n", fl.name)

	// Write content to file
	return fl.backendFile.WriteCache(ctx, data, off)
}

func (fl *file) Fsync(ctx context.Context, flags uint32) syscall.Errno {
	debugf("file[name=%q].Fsync\n", fl.name)

	// Sync cache on backend with underlying support
	return fl.backendFile.Sync(ctx)
}

func (fl *file) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	debugf("file[name=%q].Setattr\n", fl.name)
	return fl.backendFile.SetAttributes(ctx, in)
}

func (fl *file) Flush(ctx context.Context) syscall.Errno {
	debugf("file[name=%q].Flush\n", fl.name)

	// Sync cache on backend with underlying support
	return fl.backendFile.Sync(ctx)
}
