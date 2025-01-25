package implementers

import (
	"context"
	"fmt"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

// NodeImplementer is a struct that implements every callback for node
// from github.com/hanwen/go-fuse/v2/fs. It should returns an error i_ the function
// is called and should be implemented. It is used to check i_ some non implemented
// calls are called.
type NodeImplementer struct{}

// Capabilities that the FSImplementer struct should implements.
var (
	_ fs.NodeAccesser       = (*NodeImplementer)(nil)
	_ fs.NodeAllocater      = (*NodeImplementer)(nil)
	_ fs.NodeCopyFileRanger = (*NodeImplementer)(nil)
	_ fs.NodeCreater        = (*NodeImplementer)(nil)
	_ fs.NodeFlusher        = (*NodeImplementer)(nil)
	_ fs.NodeFsyncer        = (*NodeImplementer)(nil)
	_ fs.NodeGetattrer      = (*NodeImplementer)(nil)
	_ fs.NodeGetlker        = (*NodeImplementer)(nil)
	_ fs.NodeGetxattrer     = (*NodeImplementer)(nil)
	_ fs.NodeLinker         = (*NodeImplementer)(nil)
	_ fs.NodeListxattrer    = (*NodeImplementer)(nil)
	_ fs.NodeLookuper       = (*NodeImplementer)(nil)
	_ fs.NodeLseeker        = (*NodeImplementer)(nil)
	_ fs.NodeMkdirer        = (*NodeImplementer)(nil)
	_ fs.NodeMknoder        = (*NodeImplementer)(nil)
	_ fs.NodeOpendirHandler = (*NodeImplementer)(nil)
	_ fs.NodeOpendirer      = (*NodeImplementer)(nil)
	_ fs.NodeOpener         = (*NodeImplementer)(nil)
	_ fs.NodeReaddirer      = (*NodeImplementer)(nil)
	_ fs.NodeReader         = (*NodeImplementer)(nil)
	_ fs.NodeReadlinker     = (*NodeImplementer)(nil)
	_ fs.NodeReleaser       = (*NodeImplementer)(nil)
	_ fs.NodeRemovexattrer  = (*NodeImplementer)(nil)
	_ fs.NodeRenamer        = (*NodeImplementer)(nil)
	_ fs.NodeRmdirer        = (*NodeImplementer)(nil)
	_ fs.NodeSetattrer      = (*NodeImplementer)(nil)
	_ fs.NodeSetlker        = (*NodeImplementer)(nil)
	_ fs.NodeSetlkwer       = (*NodeImplementer)(nil)
	_ fs.NodeSetxattrer     = (*NodeImplementer)(nil)
	_ fs.NodeStatfser       = (*NodeImplementer)(nil)
	_ fs.NodeStatxer        = (*NodeImplementer)(nil)
	_ fs.NodeSymlinker      = (*NodeImplementer)(nil)
	_ fs.NodeUnlinker       = (*NodeImplementer)(nil)
	_ fs.NodeWriter         = (*NodeImplementer)(nil)
)

//nolint:unparam
func (ni NodeImplementer) detectorf(skippable bool, format string, args ...interface{}) {
	if skippable {
		fmt.Printf("SKIPPABLE: NodeImplementer."+format+"\n", args...)
	} else {
		fmt.Printf("NOT IMPLEMENTED: NodeImplementer."+format+"\n", args...)
	}
}

// Access is a callback of the node.
func (ni NodeImplementer) Access(_ context.Context, _ uint32) syscall.Errno {
	ni.detectorf(true, "Access")
	return fs.OK // OK i_ is not implemented
}

// Allocate is a callback of the node.
func (ni NodeImplementer) Allocate(
	_ context.Context,
	_ fs.FileHandle,
	_ uint64,
	_ uint64,
	_ uint32,
) syscall.Errno {
	ni.detectorf(false, "Allocate")
	return syscall.EOPNOTSUPP
}

// CopyFileRange is a callback of the node.
func (ni NodeImplementer) CopyFileRange(
	_ context.Context,
	_ fs.FileHandle,
	_ uint64,
	_ *fs.Inode,
	_ fs.FileHandle,
	_ uint64,
	_ uint64,
	_ uint64,
) (uint32, syscall.Errno) {
	ni.detectorf(false, "CopyFileRange")
	return 0, syscall.EOPNOTSUPP
}

// Read is a callback of the node.
func (ni NodeImplementer) Read(
	_ context.Context,
	_ fs.FileHandle,
	_ []byte,
	_ int64,
) (fuse.ReadResult, syscall.Errno) {
	ni.detectorf(false, "Read")
	return nil, syscall.EOPNOTSUPP
}

// Link is a callback of the node.
func (ni NodeImplementer) Link(
	_ context.Context,
	_ fs.InodeEmbedder,
	_ string,
	_ *fuse.EntryOut,
) (*fs.Inode, syscall.Errno) {
	ni.detectorf(false, "Link")
	return nil, syscall.EOPNOTSUPP
}

// Flush is a callback of the node.
func (ni NodeImplementer) Flush(_ context.Context, _ fs.FileHandle) syscall.Errno {
	ni.detectorf(true, "Flush")
	return fs.OK // OK i_ not implemented
}

// Fsync is a callback of the node.
func (ni NodeImplementer) Fsync(_ context.Context, _ fs.FileHandle, _ uint32) syscall.Errno {
	ni.detectorf(false, "Fsync")
	return syscall.EOPNOTSUPP
}

// Getattr is a callback of the node.
func (ni NodeImplementer) Getattr(_ context.Context, _ fs.FileHandle, _ *fuse.AttrOut) syscall.Errno {
	ni.detectorf(false, "Getattr")
	return syscall.EOPNOTSUPP
}

// Getxattr is a callback of the node.
func (ni NodeImplementer) Getxattr(_ context.Context, _ string, _ []byte) (uint32, syscall.Errno) {
	ni.detectorf(true, "Getxattr")
	return 0, syscall.EOPNOTSUPP
}

// Statfs is a callback of the node.
func (ni NodeImplementer) Statfs(_ context.Context, _ *fuse.StatfsOut) syscall.Errno {
	ni.detectorf(false, "Statfs")
	return syscall.EOPNOTSUPP
}

// Listxattr is a callback of the node.
func (ni NodeImplementer) Listxattr(_ context.Context, _ []byte) (uint32, syscall.Errno) {
	ni.detectorf(false, "Listxattr")
	return 0, syscall.EOPNOTSUPP
}

// Lseek is a callback of the node.
func (ni NodeImplementer) Lseek(
	_ context.Context,
	_ fs.FileHandle,
	_ uint64,
	_ uint32,
) (uint64, syscall.Errno) {
	ni.detectorf(false, "Lseek")
	return 0, syscall.EOPNOTSUPP
}

// Lookup is a callback of the node.
func (ni NodeImplementer) Lookup(_ context.Context, _ string, _ *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	ni.detectorf(false, "Lookup")
	return nil, syscall.EOPNOTSUPP
}

// Getlk is a callback of the node.
func (ni NodeImplementer) Getlk(
	_ context.Context,
	_ fs.FileHandle,
	_ uint64,
	_ *fuse.FileLock,
	_ uint32,
	_ *fuse.FileLock,
) syscall.Errno {
	ni.detectorf(false, "Getlk")
	return syscall.EOPNOTSUPP
}

// Mknod is a callback of the node.
func (ni NodeImplementer) Mknod(
	_ context.Context,
	_ string,
	_ uint32,
	_ uint32,
	_ *fuse.EntryOut,
) (*fs.Inode, syscall.Errno) {
	ni.detectorf(false, "Mknod")
	return nil, syscall.EOPNOTSUPP
}

// Readdir is a callback of the node.
func (ni NodeImplementer) Readdir(_ context.Context) (fs.DirStream, syscall.Errno) {
	ni.detectorf(false, "Readdir")
	return nil, syscall.EOPNOTSUPP
}

// Mkdir is a callback of the node.
func (ni NodeImplementer) Mkdir(
	_ context.Context,
	_ string,
	_ uint32,
	_ *fuse.EntryOut,
) (*fs.Inode, syscall.Errno) {
	ni.detectorf(false, "Mkdir")
	return nil, syscall.EOPNOTSUPP
}

// Create is a callback of the node.
func (ni NodeImplementer) Create(
	_ context.Context,
	_ string,
	_ uint32,
	_ uint32,
	_ *fuse.EntryOut,
) (*fs.Inode, fs.FileHandle, uint32, syscall.Errno) {
	ni.detectorf(false, "Create")
	return nil, nil, 0, syscall.EOPNOTSUPP
}

// Setxattr is a callback of the node.
func (ni NodeImplementer) Setxattr(_ context.Context, _ string, _ []byte, _ uint32) syscall.Errno {
	ni.detectorf(false, "Setxattr")
	return syscall.EOPNOTSUPP
}

// Unlink is a callback of the node.
func (ni NodeImplementer) Unlink(_ context.Context, _ string) syscall.Errno {
	ni.detectorf(false, "Unlink")
	return syscall.EOPNOTSUPP
}

// Rmdir is a callback of the node.
func (ni NodeImplementer) Rmdir(_ context.Context, _ string) syscall.Errno {
	ni.detectorf(false, "Rmdir")
	return syscall.EOPNOTSUPP
}

// Rename is a callback of the node.
func (ni NodeImplementer) Rename(
	_ context.Context,
	_ string,
	_ fs.InodeEmbedder,
	_ string,
	_ uint32,
) syscall.Errno {
	ni.detectorf(false, "Rename")
	return syscall.EOPNOTSUPP
}

// Open is a callback of the node.
func (ni NodeImplementer) Open(_ context.Context, _ uint32) (fs.FileHandle, uint32, syscall.Errno) {
	ni.detectorf(false, "Open")
	return nil, 0, syscall.EOPNOTSUPP
}

// OpendirHandle is a callback of the node.
func (ni NodeImplementer) OpendirHandle(
	_ context.Context,
	_ uint32,
) (fs.FileHandle, uint32, syscall.Errno) {
	ni.detectorf(true, "OpendirHandle")
	return nil, 0, fs.OK
}

// Opendir is a callback of the node.
func (ni NodeImplementer) Opendir(_ context.Context) syscall.Errno {
	ni.detectorf(true, "Opendir")
	return fs.OK
}

// Setattr is a callback of the node.
func (ni NodeImplementer) Setattr(
	_ context.Context,
	_ fs.FileHandle,
	_ *fuse.SetAttrIn,
	_ *fuse.AttrOut,
) syscall.Errno {
	ni.detectorf(false, "Setattr")
	return syscall.EOPNOTSUPP
}

// Readlink is a callback of the node.
func (ni NodeImplementer) Readlink(_ context.Context) ([]byte, syscall.Errno) {
	ni.detectorf(false, "Readlink")
	return nil, syscall.EOPNOTSUPP
}

// Release is a callback of the node.
func (ni NodeImplementer) Release(_ context.Context, _ fs.FileHandle) syscall.Errno {
	ni.detectorf(true, "Release")
	return syscall.EOPNOTSUPP
}

// Removexattr is a callback of the node.
func (ni NodeImplementer) Removexattr(_ context.Context, _ string) syscall.Errno {
	ni.detectorf(false, "Removexattr")
	return syscall.EOPNOTSUPP
}

// Setlk is a callback of the node.
func (ni NodeImplementer) Setlk(
	_ context.Context,
	_ fs.FileHandle,
	_ uint64,
	_ *fuse.FileLock,
	_ uint32,
) syscall.Errno {
	ni.detectorf(false, "Setlk")
	return syscall.EOPNOTSUPP
}

// Setlkw is a callback of the node.
func (ni NodeImplementer) Setlkw(
	_ context.Context,
	_ fs.FileHandle,
	_ uint64,
	_ *fuse.FileLock,
	_ uint32,
) syscall.Errno {
	ni.detectorf(false, "Setlkw")
	return syscall.EOPNOTSUPP
}

// Statx is a callback of the node.
func (ni NodeImplementer) Statx(
	_ context.Context,
	_ fs.FileHandle,
	_ uint32,
	_ uint32,
	_ *fuse.StatxOut,
) syscall.Errno {
	ni.detectorf(false, "Statx")
	return syscall.EOPNOTSUPP
}

// Symlink is a callback of the node.
func (ni NodeImplementer) Symlink(
	_ context.Context,
	_, _ string,
	_ *fuse.EntryOut,
) (*fs.Inode, syscall.Errno) {
	ni.detectorf(false, "Symlink")
	return nil, syscall.EOPNOTSUPP
}

// Write is a callback of the node.
func (ni NodeImplementer) Write(
	_ context.Context,
	_ fs.FileHandle,
	_ []byte,
	_ int64,
) (uint32, syscall.Errno) {
	ni.detectorf(false, "Write")
	return 0, syscall.EOPNOTSUPP
}
