package wrapper

import (
	"context"
	"errors"
	"io"
	"log"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/chonker"
	"golang.org/x/sys/unix"
)

type DirectoryOption func(dir *Directory)

func WithDirectoryLogger(logger *log.Logger) DirectoryOption {
	return func(dir *Directory) {
		dir.logger = logger
	}
}

func WithDirectoryChunkSize(chunkSize int) DirectoryOption {
	return func(dir *Directory) {
		dir.chunkSize = chunkSize
	}
}

// Capabilities that the dir struct should implements
var (
	_ fs.InodeEmbedder = (*Directory)(nil)

	_ fs.NodeCreater   = (*Directory)(nil)
	_ fs.NodeGetattrer = (*Directory)(nil)
	_ fs.NodeLookuper  = (*Directory)(nil)
	_ fs.NodeMkdirer   = (*Directory)(nil)
	_ fs.NodeRenamer   = (*Directory)(nil)
	_ fs.NodeReaddirer = (*Directory)(nil)
	_ fs.NodeRmdirer   = (*Directory)(nil)
	_ fs.NodeSetattrer = (*Directory)(nil)
	_ fs.NodeStatxer   = (*Directory)(nil)
	_ fs.NodeUnlinker  = (*Directory)(nil)
)

const dirMode = syscall.S_IFDIR | syscall.S_IRWXU | syscall.S_IRGRP | syscall.S_IXGRP | syscall.S_IROTH | syscall.S_IXOTH

type Directory struct {
	fs.Inode

	backend chonker.Directory

	// Optional

	options   []DirectoryOption
	logger    *log.Logger
	chunkSize int
}

func NewDirectory(backend chonker.Directory, options ...DirectoryOption) *Directory {
	// Create a default directory
	dir := &Directory{
		backend:   backend,
		options:   options,
		chunkSize: DefaultChunkSize,
		logger:    log.New(io.Discard, "", 0),
	}

	// Apply options
	for _, o := range options {
		o(dir)
	}

	return dir
}

func (d *Directory) Create(
	ctx context.Context,
	name string,
	flags uint32,
	mode uint32,
	out *fuse.EntryOut,
) (node *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	d.logger.Printf("Directory.Create(name=%q, ...)\n", name)

	// Create a new child file from backend
	backendChildFile, err := d.backend.CreateFile(ctx, name, d.chunkSize)
	if err != nil {
		return nil, nil, 0, chonker.ToSyscallErrno(err, chonker.ToSyscallErrnoOptions{
			Logger: d.logger,
		})
	}

	// Create chonkfs File
	f := NewFile(backendChildFile,
		WithFileLogger(d.logger),
		WithFileChunkSize(d.chunkSize),
		WithFileName(name))

	// Return an inode with the chonkfs directory
	return d.NewInode(ctx, f, fs.StableAttr{Mode: syscall.S_IFREG}), f, fuse.FOPEN_DIRECT_IO, fs.OK
}

func (d *Directory) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	d.logger.Printf("Directory.Getattr(...)\n")

	out.Mode = dirMode
	out.Blksize = uint32(d.chunkSize)

	return fs.OK
}

func (d *Directory) Statx(ctx context.Context, f fs.FileHandle, flags uint32, mask uint32, out *fuse.StatxOut) syscall.Errno {
	d.logger.Printf("Directory.Statx(...)\n")

	out.Mode = dirMode
	out.Blksize = uint32(d.chunkSize)

	return fs.OK
}

func (d *Directory) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	d.logger.Printf("Directory.Lookup(name=%q, ...)\n", name)

	// Get backend child directory
	backendChildDir, err := d.backend.GetDirectory(ctx, name)

	switch err {
	case nil:
		// Create inode
		ino := d.NewInode(ctx,
			NewDirectory(backendChildDir, d.options...),
			fs.StableAttr{
				Mode: syscall.S_IFDIR,
			})

		// Set mode from backend
		_, err := backendChildDir.GetAttributes(ctx)
		if err != nil {
			return nil, chonker.ToSyscallErrno(err, chonker.ToSyscallErrnoOptions{
				Logger: d.logger,
			})
		}

		// Add info
		out.Blksize = uint32(d.chunkSize)
		out.Mode = dirMode

		// Return the inode
		return ino, fs.OK
	case chonker.ErrNotDirectory:
		// Get backend file
		backendChildFile, err := d.backend.GetFile(ctx, name)
		if err != nil {
			return nil, chonker.ToSyscallErrno(err, chonker.ToSyscallErrnoOptions{
				Logger: d.logger,
			})
		}

		// Create inode
		ino := d.NewInode(ctx,
			NewFile(backendChildFile,
				WithFileLogger(d.logger),
				WithFileChunkSize(d.chunkSize),
				WithFileName(name)),
			fs.StableAttr{
				Mode: syscall.S_IFREG,
			})

		// Set mode from backend
		attr, err := backendChildFile.GetAttributes(ctx)
		if err != nil {
			return nil, chonker.ToSyscallErrno(err, chonker.ToSyscallErrnoOptions{
				Logger: d.logger,
			})
		}

		// Add info
		out.Size = uint64(attr.Size)
		out.Blocks = uint64((attr.Size-1)/d.chunkSize + 1)
		out.Blksize = uint32(d.chunkSize)
		out.Mode = fileMode

		// Return the inode
		return ino, fs.OK
	default:
		return nil, chonker.ToSyscallErrno(err, chonker.ToSyscallErrnoOptions{
			Logger: d.logger,
		})
	}
}

func (d *Directory) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	d.logger.Printf("Directory.Mkdir(...)\n")

	// Create a new child directory from backend
	backendChildDir, err := d.backend.CreateDirectory(ctx, name)
	if err != nil {
		return nil, chonker.ToSyscallErrno(err, chonker.ToSyscallErrnoOptions{
			Logger: d.logger,
		})
	}

	// Return an inode with the chonkfs directory
	return d.NewInode(ctx,
		NewDirectory(backendChildDir, d.options...),
		fs.StableAttr{Mode: syscall.S_IFDIR}), fs.OK
}

func (d *Directory) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	d.logger.Printf("Directory.Readdir(...)\n")

	list := make([]fuse.DirEntry, 0)

	// Get directories
	dirs, err := d.backend.ListDirectories(ctx)
	if err != nil {
		return nil, chonker.ToSyscallErrno(err, chonker.ToSyscallErrnoOptions{
			Logger: d.logger,
		})
	}

	// Add directories
	for _, name := range dirs {
		list = append(list, fuse.DirEntry{
			Name: name,
			Mode: fuse.S_IFDIR,
		})
	}

	// Get files
	files, err := d.backend.ListFiles(ctx)
	if err != nil {
		return nil, chonker.ToSyscallErrno(err, chonker.ToSyscallErrnoOptions{
			Logger: d.logger,
		})
	}

	// Add files
	for _, name := range files {
		list = append(list, fuse.DirEntry{
			Name: name,
			Mode: fuse.S_IFREG,
		})
	}

	return fs.NewListDirStream(list), fs.OK
}

func (d *Directory) Rmdir(ctx context.Context, name string) syscall.Errno {
	d.logger.Printf("Directory.Rmdir(...)\n")
	return chonker.ToSyscallErrno(
		d.backend.RemoveDirectory(ctx, name),
		chonker.ToSyscallErrnoOptions{
			Logger: d.logger,
		},
	)
}

func (d *Directory) Unlink(ctx context.Context, name string) syscall.Errno {
	d.logger.Printf("Directory.Unlink(name=%q, ...)\n", name)
	return chonker.ToSyscallErrno(
		d.backend.RemoveFile(ctx, name),
		chonker.ToSyscallErrnoOptions{
			Logger: d.logger,
		},
	)
}

func (d *Directory) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	d.logger.Printf("Directory.Setattr(...)\n")
	return fs.OK
}

func (d *Directory) Rename(ctx context.Context, name string, newParent fs.InodeEmbedder, newName string, flags uint32) syscall.Errno {
	d.logger.Printf("Directory.Rename(name=%q, newName=%q)\n", name, newName)

	// Get the new parent directory
	newParentDir, errno := d.getDirectoryFromInodeEmbedder(newParent)
	if errno != fs.OK {
		return errno
	}

	// Check if no replace flag is set
	noReplace := (flags & unix.RENAME_SECLUDE) == unix.RENAME_SECLUDE

	// Check if the directory or file exists
	_, err := d.backend.GetDirectory(ctx, name)

	switch {
	case err != nil && !errors.Is(err, chonker.ErrNotDirectory):
		// Error happened
		return chonker.ToSyscallErrno(err, chonker.ToSyscallErrnoOptions{
			Logger: d.logger,
		})
	case errors.Is(err, chonker.ErrNotDirectory):
		// It's a file
		if err := d.backend.RenameFile(ctx, name, newParentDir.backend, newName, noReplace); err != nil {
			return chonker.ToSyscallErrno(err, chonker.ToSyscallErrnoOptions{
				Logger: d.logger,
			})
		}
	default:
		// No error: its a directory
		if err := d.backend.RenameDirectory(ctx, name, newParentDir.backend, newName, noReplace); err != nil {
			return chonker.ToSyscallErrno(err, chonker.ToSyscallErrnoOptions{
				Logger: d.logger,
			})
		}
	}

	return fs.OK
}

func (d *Directory) getDirectoryFromInodeEmbedder(inode fs.InodeEmbedder) (*Directory, syscall.Errno) {
	// Cast/Assert new parent to directory structure
	dir, ok := inode.(*Directory)
	if !ok {
		log.Printf("ERROR: new parent is not a ChonkFS directory\n")
		return nil, syscall.EINVAL
	}

	return dir, fs.OK
}
