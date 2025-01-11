package implementers

import (
	"context"
	"fmt"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

// NodeImplementer is a struct that implements every callback for node
// from github.com/hanwen/go-fuse/v2/fs. It should returns an error if the function
// is called and should be implemented. It is used to check if some non implemented
// calls are called.
type NodeImplementer struct{}

// Capabilities that the FSImplementer struct should implements
var (
	_ fs.NodeAccesser       = (*NodeImplementer)(nil)
	_ fs.NodeAllocater      = (*NodeImplementer)(nil)
	_ fs.NodeCopyFileRanger = (*NodeImplementer)(nil)
	_ fs.NodeCreater        = (*NodeImplementer)(nil)
	_ fs.NodeFlusher        = (*NodeImplementer)(nil)
	_ fs.NodeFsyncer        = (*NodeImplementer)(nil)
	_ fs.NodeGetattrer      = (*NodeImplementer)(nil)
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
	_ fs.NodeReadlinker     = (*NodeImplementer)(nil)
	_ fs.NodeReleaser       = (*NodeImplementer)(nil)
	_ fs.NodeRemovexattrer  = (*NodeImplementer)(nil)
	_ fs.NodeRenamer        = (*NodeImplementer)(nil)
	_ fs.NodeRmdirer        = (*NodeImplementer)(nil)
	_ fs.NodeSetattrer      = (*NodeImplementer)(nil)
	_ fs.NodeSetlker        = (*NodeImplementer)(nil)
	_ fs.NodeSetlkwer       = (*NodeImplementer)(nil)
	_ fs.NodeStatfser       = (*NodeImplementer)(nil)
	_ fs.NodeStatxer        = (*NodeImplementer)(nil)
	_ fs.NodeSymlinker      = (*NodeImplementer)(nil)
	_ fs.NodeUnlinker       = (*NodeImplementer)(nil)
	_ fs.NodeWriter         = (*NodeImplementer)(nil)
)

func (ni NodeImplementer) Detector(skippable bool, format string, args ...interface{}) {
	if skippable {
		fmt.Printf("SKIPPABLE: NodeImplementer."+format, args...)
	} else {
		fmt.Printf("NOT IMPLEMENTED: NodeImplementer."+format, args...)
	}
}

func (ni NodeImplementer) Access(ctx context.Context, mask uint32) syscall.Errno {
	ni.Detector(true, "Access\n")
	return fs.OK // OK if is not implemented
}

func (ni NodeImplementer) Allocate(ctx context.Context, f fs.FileHandle, off uint64, size uint64, mode uint32) syscall.Errno {
	ni.Detector(false, "Allocate\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) CopyFileRange(
	ctx context.Context,
	fhIn fs.FileHandle,
	offIn uint64,
	out *fs.Inode,
	fhOut fs.FileHandle,
	offOut uint64,
	len uint64,
	flags uint64,
) (uint32, syscall.Errno) {
	ni.Detector(false, "CopyFileRange\n")
	return 0, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Link(
	ctx context.Context,
	target fs.InodeEmbedder,
	ame string,
	out *fuse.EntryOut,
) (node *fs.Inode, errno syscall.Errno) {
	ni.Detector(false, "Link\n")
	return nil, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Flush(ctx context.Context, f fs.FileHandle) syscall.Errno {
	ni.Detector(true, "Flush\n")
	return fs.OK // OK if not implemented
}

func (ni NodeImplementer) Fsync(ctx context.Context, f fs.FileHandle, flags uint32) syscall.Errno {
	ni.Detector(false, "Fsync\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	ni.Detector(false, "Getattr\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Statfs(ctx context.Context, out *fuse.StatfsOut) syscall.Errno {
	ni.Detector(false, "Statfs\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Listxattr(ctx context.Context, dest []byte) (uint32, syscall.Errno) {
	ni.Detector(false, "Listxattr\n")
	return 0, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Lseek(ctx context.Context, f fs.FileHandle, Off uint64, whence uint32) (uint64, syscall.Errno) {
	ni.Detector(false, "Lseek\n")
	return 0, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	ni.Detector(false, "Lookup\n")
	return nil, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Mknod(
	ctx context.Context,
	name string,
	mode uint32,
	dev uint32,
	out *fuse.EntryOut,
) (*fs.Inode, syscall.Errno) {
	ni.Detector(false, "Mknod\n")
	return nil, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	ni.Detector(false, "Readdir\n")
	return nil, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	ni.Detector(false, "Mkdir\n")
	return nil, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Create(
	ctx context.Context,
	name string,
	flags uint32,
	mode uint32,
	out *fuse.EntryOut,
) (node *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	ni.Detector(false, "Create\n")
	return nil, nil, 0, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Unlink(ctx context.Context, name string) syscall.Errno {
	ni.Detector(false, "Unlink\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Rmdir(ctx context.Context, name string) syscall.Errno {
	ni.Detector(false, "Rmdir\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Rename(ctx context.Context, name string, newParent fs.InodeEmbedder, newName string, flags uint32) syscall.Errno {
	ni.Detector(false, "Rename\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	ni.Detector(false, "Open\n")
	return nil, 0, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) OpendirHandle(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	ni.Detector(true, "OpendirHandle\n")
	return nil, 0, fs.OK
}

func (ni NodeImplementer) Opendir(ctx context.Context) syscall.Errno {
	ni.Detector(true, "Opendir\n")
	return fs.OK
}

func (ni NodeImplementer) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	ni.Detector(false, "Setattr\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Readlink(ctx context.Context) ([]byte, syscall.Errno) {
	ni.Detector(false, "Readlink\n")
	return nil, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Release(ctx context.Context, f fs.FileHandle) syscall.Errno {
	ni.Detector(true, "Release\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Removexattr(ctx context.Context, attr string) syscall.Errno {
	ni.Detector(false, "Removexattr\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Setlk(ctx context.Context, f fs.FileHandle, owner uint64, lk *fuse.FileLock, flags uint32) syscall.Errno {
	ni.Detector(false, "Setlk\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Setlkw(ctx context.Context, f fs.FileHandle, owner uint64, lk *fuse.FileLock, flags uint32) syscall.Errno {
	ni.Detector(false, "Setlkw\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Statx(ctx context.Context, f fs.FileHandle, flags uint32, mask uint32, out *fuse.StatxOut) syscall.Errno {
	ni.Detector(false, "Statx\n")
	return syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Symlink(ctx context.Context, target, name string, out *fuse.EntryOut) (node *fs.Inode, errno syscall.Errno) {
	ni.Detector(false, "Symlink\n")
	return nil, syscall.EOPNOTSUPP
}

func (ni NodeImplementer) Write(ctx context.Context, f fs.FileHandle, data []byte, off int64) (written uint32, errno syscall.Errno) {
	ni.Detector(false, "Write\n")
	return 0, syscall.EOPNOTSUPP
}
