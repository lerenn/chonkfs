package chonkfs

import (
	"context"
	"io"
	"log"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/backends"
)

type FileOption func(fl *File)

func WithFileLogger(logger *log.Logger) FileOption {
	return func(fl *File) {
		fl.logger = logger
	}
}

func WithFileChunkSize(chunkSize int) FileOption {
	return func(fl *File) {
		fl.chunkSize = chunkSize
	}
}

func WithFileName(name string) FileOption {
	return func(fl *File) {
		fl.name = name
	}
}

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
	fs.Inode

	backend      backends.File
	sessionFlags uint32

	// Optional

	options   []FileOption
	logger    *log.Logger
	name      string
	chunkSize int
}

func NewFile(backend backends.File, options ...FileOption) *File {
	// Create default file
	f := &File{
		backend:   backend,
		options:   options,
		logger:    log.New(io.Discard, "", 0),
		name:      "",
		chunkSize: DefaultChunkSize,
	}

	// Execute options
	for _, o := range options {
		o(f)
	}

	return f
}

func (fl *File) Getattr(ctx context.Context, out *fuse.AttrOut) (errno syscall.Errno) {
	fl.logger.Printf("File[%s].Getattr(...)\n", fl.name)

	// Get attributes from backend
	attr, err := fl.backend.GetAttributes(ctx)
	if err != nil {
		return backends.ToSyscallErrno(err,
			backends.ToSyscallErrnoOptions{
				Logger: fl.logger,
			})
	}

	// Set attributes
	out.Attr = attr

	return fs.OK
}

func (fl *File) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	fl.logger.Printf("File[%s].Read(len=%d, off=%d)\n", fl.name, len(dest), off)

	// Get content from file
	start := uint64(off)
	end := uint64(off) + uint64(len(dest))
	content, err := fl.backend.Read(ctx, int(start), int(end))
	if err != nil {
		return nil, backends.ToSyscallErrno(err,
			backends.ToSyscallErrnoOptions{
				Logger: fl.logger,
			})
	}

	return fuse.ReadResultData(content), fs.OK
}

func (fl *File) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	fl.logger.Printf("File[%s].Open(...)\n", fl.name)

	// Save flags
	fl.sessionFlags = flags

	// Check if file exists if O_EXCL
	// TODO

	return fl, fuse.FOPEN_DIRECT_IO, fs.OK
}

func (fl *File) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	fl.logger.Printf("File[%s].Write(len=%d, off=%d)\n", fl.name, len(data), off)

	// Write content to file
	w, err := fl.backend.Write(ctx, data, int(off), backends.WriteOptions{
		Truncate: fl.sessionFlags&syscall.O_TRUNC != 0,
		Append:   fl.sessionFlags&syscall.O_APPEND != 0,
	})
	return uint32(w), backends.ToSyscallErrno(err,
		backends.ToSyscallErrnoOptions{
			Logger: fl.logger,
		})
}

func (fl *File) Fsync(ctx context.Context, flags uint32) syscall.Errno {
	fl.logger.Printf("File[%s].Fsync(...)\n", fl.name)

	// Sync cache on backend with underlying support
	return backends.ToSyscallErrno(
		fl.backend.Sync(ctx),
		backends.ToSyscallErrnoOptions{
			Logger: fl.logger,
		})
}

func (fl *File) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	fl.logger.Printf("File[%s].Setattr(in=%+v, out=%+v)\n", fl.name, *in, *out)

	// Get actual size
	actualSize, err := fl.backend.Size(ctx)
	if err != nil {
		return backends.ToSyscallErrno(err,
			backends.ToSyscallErrnoOptions{
				Logger: fl.logger,
			})
	}

	// Truncate the file if needed
	if in.Size < uint64(actualSize) {
		if err := fl.backend.Truncate(ctx, int(in.Size)); err != nil {
			return backends.ToSyscallErrno(err,
				backends.ToSyscallErrnoOptions{
					Logger: fl.logger,
				})
		}
	}

	return backends.ToSyscallErrno(
		fl.backend.SetAttributes(ctx, in),
		backends.ToSyscallErrnoOptions{
			Logger: fl.logger,
		})
}

func (fl *File) Flush(ctx context.Context) syscall.Errno {
	fl.logger.Printf("File[%s].Flush(...)\n", fl.name)

	// Sync cache on backend with underlying support
	return backends.ToSyscallErrno(
		fl.backend.Sync(ctx),
		backends.ToSyscallErrnoOptions{
			Logger: fl.logger,
		})
}
