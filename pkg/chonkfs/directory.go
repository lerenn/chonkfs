package chonkfs

import (
	"context"
	"log"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/backends"
)

type Directory struct {
	backend backends.Directory
	logger  *log.Logger

	fs.Inode

	// implementers.NodeImplementer
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
	_ fs.NodeUnlinker  = (*Directory)(nil)
)

func (d *Directory) Create(
	ctx context.Context,
	name string,
	flags uint32,
	mode uint32,
	out *fuse.EntryOut,
) (node *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	d.logger.Printf("Directory.Create(name=%q, ...)\n", name)

	// Create a new child file from backend
	backendChildFile, errno := d.backend.CreateFile(ctx, name)
	if errno != fs.OK {
		return nil, nil, 0, errno
	}

	// Create chonkfs File
	f := &File{
		backend: backendChildFile,
		logger:  d.logger,
		name:    name,
	}

	// Return an inode with the chonkfs directory
	return d.NewInode(ctx, f, fs.StableAttr{Mode: syscall.S_IFREG}), f, fuse.FOPEN_DIRECT_IO, fs.OK
}

func (d *Directory) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	d.logger.Printf("Directory.Getattr(...)\n")

	// Get attributes from backend
	attr, errno := d.backend.GetAttributes(ctx)
	if errno != fs.OK {
		return errno
	}

	// Set attributes
	out.Attr = attr

	return fs.OK
}

func (d *Directory) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	d.logger.Printf("Directory.Lookup(name=%q, ...)\n", name)

	// Get backend child directory
	backendChildDir, errno := d.backend.GetDirectory(ctx, name)

	switch errno {
	case fs.OK:
		// Create inode
		ino := d.NewInode(ctx, &Directory{
			backend: backendChildDir,
		}, fs.StableAttr{Mode: syscall.S_IFDIR})

		// Set mode from backend
		attr, errnoBackend := backendChildDir.GetAttributes(ctx)
		if errnoBackend != fs.OK {
			return nil, errnoBackend
		}
		out.Attr = attr

		// Add some info
		out.Mode = 1755 // TODO: fixme
		out.Ino = ino.StableAttr().Ino

		// Return the inode
		return ino, fs.OK
	case syscall.ENOTDIR:
		// Get backend file
		backendChildFile, errno := d.backend.GetFile(ctx, name)
		if errno != fs.OK {
			return nil, errno
		}

		// Create inode
		ino := d.NewInode(ctx, &File{
			backend: backendChildFile,
			logger:  d.logger,
			name:    name,
		}, fs.StableAttr{Mode: syscall.S_IFREG})

		// Set mode from backend
		attr, errnoBackend := backendChildFile.GetAttributes(ctx)
		if errnoBackend != fs.OK {
			return nil, errnoBackend
		}
		out.Attr = attr

		// Add some info
		out.Mode = 0755 // TODO: fixme
		out.Ino = ino.StableAttr().Ino

		// Return the inode
		return ino, fs.OK
	default:
		return nil, errno
	}
}

func (d *Directory) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	d.logger.Printf("Directory.Mkdir(...)\n")

	// Create a new child directory from backend
	backendChildDir, errno := d.backend.CreateDirectory(ctx, name)
	if errno != fs.OK {
		return nil, errno
	}

	// Return an inode with the chonkfs directory
	return d.NewInode(ctx, &Directory{
		backend: backendChildDir,
	}, fs.StableAttr{Mode: syscall.S_IFDIR}), fs.OK
}

func (d *Directory) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	d.logger.Printf("Directory.Readdir(...)\n")

	// List entries from backend
	l, errno := d.backend.ListEntries(ctx)
	if errno != fs.OK {
		return nil, errno
	}

	return fs.NewListDirStream(l), fs.OK
}

func (d *Directory) Rmdir(ctx context.Context, name string) syscall.Errno {
	d.logger.Printf("Directory.Rmdir(...)\n")
	return d.backend.RemoveDirectory(ctx, name)
}

func (d *Directory) Unlink(ctx context.Context, name string) syscall.Errno {
	d.logger.Printf("Directory.Unlink(name=%q, ...)\n", name)
	return d.backend.RemoveFile(ctx, name)
}

func (d *Directory) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	d.logger.Printf("Directory.Setattr(...)\n")
	return d.backend.SetAttributes(ctx, in)
}

func (d *Directory) Rename(ctx context.Context, name string, newParent fs.InodeEmbedder, newName string, flags uint32) syscall.Errno {
	d.logger.Printf("Directory.Rename(...)\n")

	// Get the new parent directory
	newParentDir, errno := getDirectoryFromInodeEmbedder(newParent)
	if errno != fs.OK {
		return errno
	}

	// Rename node on backend
	if errno := d.backend.RenameNode(ctx, name, newParentDir.backend, newName); errno != fs.OK {
		return errno
	}

	return fs.OK
}

func getDirectoryFromInodeEmbedder(inode fs.InodeEmbedder) (*Directory, syscall.Errno) {
	// Cast/Assert new parent to directory structure
	dir, ok := inode.(*Directory)
	if !ok {
		return nil, syscall.EINVAL
	}

	return dir, fs.OK
}
