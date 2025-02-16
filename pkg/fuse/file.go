package fuse

import (
	"context"
	"io"
	"log"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/chonker"
)

type fileOption func(fl *File)

// WithFileLogger is an option to set the logger of a file.
//
//nolint:revive
func WithFileLogger(logger *log.Logger) fileOption {
	return func(fl *File) {
		fl.logger = logger
	}
}

// WithFileChunkSize is an option to set the chunk size of a file.
//
//nolint:revive
func WithFileChunkSize(chunkSize int) fileOption {
	return func(fl *File) {
		fl.chunkSize = chunkSize
	}
}

// WithFileName is an option to set the name of a file.
//
//nolint:revive
func WithFileName(name string) fileOption {
	return func(fl *File) {
		fl.name = name
	}
}

// Capabilities that the file struct should implements.
var (
	_ fs.FileFlusher   = (*File)(nil)
	_ fs.FileGetattrer = (*File)(nil)
	_ fs.FileReader    = (*File)(nil)
	_ fs.FileWriter    = (*File)(nil)
	_ fs.FileFsyncer   = (*File)(nil)
	_ fs.FileStatxer   = (*File)(nil)

	_ fs.InodeEmbedder = (*File)(nil)

	_ fs.NodeOpener    = (*File)(nil)
	_ fs.NodeSetattrer = (*File)(nil)
)

const fileMode = syscall.S_IFREG | syscall.S_IRWXU | syscall.S_IRGRP |
	syscall.S_IXGRP | syscall.S_IROTH | syscall.S_IXOTH

// File is a representation of a FUSE file as wrapper of chonker.
type File struct {
	fs.Inode

	backend      chonker.File
	sessionFlags uint32

	// Optional

	options   []fileOption
	logger    *log.Logger
	name      string
	chunkSize int
}

// PreHook is a hook that is called before the file is used.
func (f *File) PreHook() {}

// PostHook is a hook that is called after the file is used.
func (f *File) PostHook() {}

// NewFile creates a new file.
func NewFile(backend chonker.File, options ...fileOption) *File {
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

// Getattr returns the attributes of the file to the FUSE system.
func (f *File) Getattr(ctx context.Context, out *fuse.AttrOut) (errno syscall.Errno) {
	f.PreHook()
	defer f.PostHook()
	f.logger.Printf("File[%s].Getattr(...)\n", f.name)

	// Get attributes from backend
	attr, err := f.backend.GetAttributes(ctx)
	if err != nil {
		return chonker.ToSyscallErrno(err,
			chonker.ToSyscallErrnoOptions{
				Logger: f.logger,
			})
	}

	// Set attributes
	out.Mode = fileMode
	out.Size = uint64(attr.Size)
	out.Blocks = uint64((attr.Size-1)/f.chunkSize + 1)
	out.Blksize = uint32(f.chunkSize)

	return fs.OK
}

// Statx returns the attributes of the file to the FUSE system.
func (f *File) Statx(ctx context.Context, _ uint32, _ uint32, out *fuse.StatxOut) syscall.Errno {
	f.PreHook()
	defer f.PostHook()
	f.logger.Printf("File[%s].Statx(...)\n", f.name)

	// Get attributes from backend
	attr, err := f.backend.GetAttributes(ctx)
	if err != nil {
		return chonker.ToSyscallErrno(err,
			chonker.ToSyscallErrnoOptions{
				Logger: f.logger,
			})
	}

	// Set attributes
	out.Mode = fileMode
	out.Size = uint64(attr.Size)
	out.Blocks = uint64((attr.Size-1)/f.chunkSize + 1)
	out.Blksize = uint32(f.chunkSize)

	return fs.OK
}

// Read reads the file for the FUSE system.
func (f *File) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	f.PreHook()
	defer f.PostHook()
	f.logger.Printf("File[%s].Read(len=%d, off=%d)\n", f.name, len(dest), off)

	// Get content from file
	dest, err := f.backend.Read(ctx, dest, int(off))
	if err != nil {
		return nil, chonker.ToSyscallErrno(err,
			chonker.ToSyscallErrnoOptions{
				Logger: f.logger,
			})
	}

	return fuse.ReadResultData(dest), fs.OK
}

// Open opens the file for the FUSE system.
func (f *File) Open(_ context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	f.PreHook()
	defer f.PostHook()
	f.logger.Printf("File[%s].Open(...)\n", f.name)

	// Save flags
	f.sessionFlags = flags

	// Check if file exists if O_EXCL
	// TODO

	return f, fuse.FOPEN_DIRECT_IO, fs.OK
}

// Write writes the file for the FUSE system.
func (f *File) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	f.PreHook()
	defer f.PostHook()
	f.logger.Printf("File[%s].Write(len=%d, off=%d)\n", f.name, len(data), off)

	// Write content to file
	w, err := f.backend.Write(ctx, data, int(off), chonker.WriteOptions{
		Truncate: f.sessionFlags&syscall.O_TRUNC != 0,
		Append:   f.sessionFlags&syscall.O_APPEND != 0,
	})
	return uint32(w), chonker.ToSyscallErrno(err,
		chonker.ToSyscallErrnoOptions{
			Logger: f.logger,
		})
}

// Fsync flushes the file for the FUSE system.
func (f *File) Fsync(ctx context.Context, _ uint32) syscall.Errno {
	f.PreHook()
	defer f.PostHook()
	f.logger.Printf("File[%s].Fsync(...)\n", f.name)

	// Sync cache on backend with underlying support
	return chonker.ToSyscallErrno(
		f.backend.Sync(ctx),
		chonker.ToSyscallErrnoOptions{
			Logger: f.logger,
		})
}

// Setattr sets the attributes of the file for the FUSE system.
func (f *File) Setattr(ctx context.Context, _ fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	f.PreHook()
	defer f.PostHook()
	f.logger.Printf("File[%s].Setattr(in=%+v, out=%+v)\n", f.name, *in, *out)

	// Get actual size
	info, err := f.backend.GetAttributes(ctx)
	if err != nil {
		return chonker.ToSyscallErrno(err,
			chonker.ToSyscallErrnoOptions{
				Logger: f.logger,
			})
	}
	actualSize := info.Size

	// Truncate the file if needed
	if in.Size < uint64(actualSize) {
		if err := f.backend.Truncate(ctx, int(in.Size)); err != nil {
			return chonker.ToSyscallErrno(err,
				chonker.ToSyscallErrnoOptions{
					Logger: f.logger,
				})
		}
	}

	return chonker.ToSyscallErrno(
		f.backend.SetAttributes(ctx, chonker.FileAttributes{}),
		chonker.ToSyscallErrnoOptions{
			Logger: f.logger,
		})
}

// Flush flushes the file for the FUSE system.
func (f *File) Flush(ctx context.Context) syscall.Errno {
	f.PreHook()
	defer f.PostHook()
	f.logger.Printf("File[%s].Flush(...)\n", f.name)

	// Sync cache on backend with underlying support
	return chonker.ToSyscallErrno(
		f.backend.Sync(ctx),
		chonker.ToSyscallErrnoOptions{
			Logger: f.logger,
		})
}
