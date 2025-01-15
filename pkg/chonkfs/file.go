package chonkfs

import (
	"context"
	"log"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/backends"
)

// Capabilities that the file struct should implements
var (
	_ fs.FileFlusher   = (*File)(nil)
	_ fs.FileGetattrer = (*File)(nil)
	_ fs.FileReader    = (*File)(nil)
	_ fs.FileWriter    = (*File)(nil)
	_ fs.FileFsyncer   = (*File)(nil)

	_ fs.InodeEmbedder = (*File)(nil)

	_ fs.NodeOpener    = (*File)(nil)
	_ fs.NodeSetattrer = (*File)(nil)
)

type File struct {
	backend backends.File
	logger  *log.Logger
	name    string

	fs.Inode

	// implementers.NodeImplementer
	// implementers.FileImplementer
}

func (fl *File) Getattr(ctx context.Context, out *fuse.AttrOut) (errno syscall.Errno) {
	fl.logger.Printf("File[%s].Getattr(...)\n", fl.name)

	// Get attributes from backend
	attr, errno := fl.backend.GetAttributes(ctx)
	if errno != fs.OK {
		return errno
	}

	// Set attributes
	out.Attr = attr

	return fs.OK
}

func (fl *File) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	fl.logger.Printf("File[%s].Read(...)\n", fl.name)

	// Get content from file
	content, errno := fl.backend.Read(ctx, off)
	if errno != fs.OK {
		return nil, errno
	}

	return fuse.ReadResultData(content), fs.OK
}

func (fl *File) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	fl.logger.Printf("File[%s].Open(...)\n", fl.name)

	// Nothing to do for the moment.
	// Please open a ticket if needed.

	return fl, fuse.FOPEN_DIRECT_IO, fs.OK
}

func (fl *File) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	fl.logger.Printf("File[%s].Write(...)\n", fl.name)

	// Write content to file
	return fl.backend.WriteCache(ctx, data, off)
}

func (fl *File) Fsync(ctx context.Context, flags uint32) syscall.Errno {
	fl.logger.Printf("File[%s].Fsync(...)\n", fl.name)

	// Sync cache on backend with underlying support
	return fl.backend.Sync(ctx)
}

func (fl *File) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	fl.logger.Printf("File[%s].Setattr(...)\n", fl.name)
	return fl.backend.SetAttributes(ctx, in)
}

func (fl *File) Flush(ctx context.Context) syscall.Errno {
	fl.logger.Printf("File[%s].Flush(...)\n", fl.name)

	// Sync cache on backend with underlying support
	return fl.backend.Sync(ctx)
}
