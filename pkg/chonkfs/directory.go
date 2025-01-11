package chonkfs

import (
	"context"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/backend"
)

type directory struct {
	backendDirectory backend.Directory

	fs.Inode

	// implementers.NodeImplementer
}

// Capabilities that the dir struct should implements
var (
	_ fs.InodeEmbedder = (*directory)(nil)

	_ fs.NodeCreater   = (*directory)(nil)
	_ fs.NodeGetattrer = (*directory)(nil)
	_ fs.NodeLookuper  = (*directory)(nil)
	_ fs.NodeMkdirer   = (*directory)(nil)
	_ fs.NodeReaddirer = (*directory)(nil)
)

func (d *directory) Create(
	ctx context.Context,
	name string,
	flags uint32,
	mode uint32,
	out *fuse.EntryOut,
) (node *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	debugf("directory.Create\n")

	// Create a new child file from backend
	backendChildFile, errno := d.backendDirectory.CreateChildFile(ctx, name)
	if errno != fs.OK {
		return nil, nil, 0, errno
	}

	// Return an inode with the chonkfs directory
	// TODO: implement file handler
	return d.NewInode(ctx, &file{
		backendFile: backendChildFile,
	}, fs.StableAttr{Mode: syscall.S_IFREG}), 0, 0, fs.OK
}

func (d *directory) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	debugf("directory.Getattr\n")

	// Nothing to do for the moment.
	// Please open a ticket if needed.

	return fs.OK
}

func (d *directory) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	debugf("directory.Lookup\n")

	// Get backend child directory
	backendChildDir, errno := d.backendDirectory.GetChildDirectory(ctx, name)

	switch errno {
	case fs.OK:
		// Set mode
		// TODO: fixme
		out.Mode = 1755

		// Return an inode with the chonkfs directory
		return d.NewInode(ctx, &directory{
			backendDirectory: backendChildDir,
		}, fs.StableAttr{Mode: syscall.S_IFDIR}), fs.OK
	case syscall.ENOTDIR:
		// Get backend file
		backendChildFile, errno := d.backendDirectory.GetChildFile(ctx, name)
		if errno != fs.OK {
			return nil, errno
		}

		// Set mode
		// TODO: fixme
		out.Mode = 0755

		// Return an inode with the chonkfs directory
		return d.NewInode(ctx, &file{
			backendFile: backendChildFile,
		}, fs.StableAttr{Mode: syscall.S_IFREG}), fs.OK

	default:
		return nil, errno
	}

}

func (d *directory) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	debugf("directory.Mkdir\n")

	// Create a new child directory from backend
	backendChildDir, errno := d.backendDirectory.CreateChildDirectory(ctx, name)
	if errno != fs.OK {
		return nil, errno
	}

	// Return an inode with the chonkfs directory
	return d.NewInode(ctx, &directory{
		backendDirectory: backendChildDir,
	}, fs.StableAttr{Mode: syscall.S_IFDIR}), fs.OK
}

func (d *directory) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	debugf("directory.Readdir\n")

	// List entries from backend
	l, errno := d.backendDirectory.ListDirectoryEntries(ctx)
	if errno != fs.OK {
		return nil, errno
	}

	return fs.NewListDirStream(l), fs.OK
}
