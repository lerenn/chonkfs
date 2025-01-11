package implementers

import (
	"context"
	"fmt"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

// FileImplementer is a struct that implements every callback for file
// from github.com/hanwen/go-fuse/v2/fs. It should returns an error if the function
// is called and should be implemented. It is used to check if some non implemented
// calls are called.
type FileImplementer struct{}

// Capabilities that the FSImplementer struct should implements
var (
	_ fs.FileAllocater       = (*FileImplementer)(nil)
	_ fs.FileFlusher         = (*FileImplementer)(nil)
	_ fs.FileFsyncdirer      = (*FileImplementer)(nil)
	_ fs.FileFsyncer         = (*FileImplementer)(nil)
	_ fs.FileGetattrer       = (*FileImplementer)(nil)
	_ fs.FileGetlker         = (*FileImplementer)(nil)
	_ fs.FileLseeker         = (*FileImplementer)(nil)
	_ fs.FilePassthroughFder = (*FileImplementer)(nil)
	_ fs.FileReaddirenter    = (*FileImplementer)(nil)
	_ fs.FileReader          = (*FileImplementer)(nil)
	_ fs.FileReleasedirer    = (*FileImplementer)(nil)
	_ fs.FileReleaser        = (*FileImplementer)(nil)
	_ fs.FileSeekdirer       = (*FileImplementer)(nil)
	_ fs.FileSetattrer       = (*FileImplementer)(nil)
	_ fs.FileSetlker         = (*FileImplementer)(nil)
	_ fs.FileSetlkwer        = (*FileImplementer)(nil)
	_ fs.FileStatxer         = (*FileImplementer)(nil)
	_ fs.FileWriter          = (*FileImplementer)(nil)
)

func (fi FileImplementer) Detector(skippable bool, format string, args ...interface{}) {
	if skippable {
		fmt.Printf("SKIPPABLE: FileImplementer."+format+"\n", args...)
	} else {
		fmt.Printf("NOT IMPLEMENTED: FileImplementer."+format+"\n", args...)
	}
}

func (fi FileImplementer) Allocate(ctx context.Context, off uint64, size uint64, mode uint32) syscall.Errno {
	fi.Detector(false, "Allocate")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Flush(ctx context.Context) syscall.Errno {
	fi.Detector(false, "Flush")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Fsyncdir(ctx context.Context, flags uint32) syscall.Errno {
	fi.Detector(false, "Fsyncdir")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Fsync(ctx context.Context, flags uint32) syscall.Errno {
	fi.Detector(false, "Fsync")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Getattr(ctx context.Context, out *fuse.AttrOut) syscall.Errno {
	fi.Detector(false, "Getattr")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Getlk(ctx context.Context, owner uint64, lk *fuse.FileLock, flags uint32, out *fuse.FileLock) syscall.Errno {
	fi.Detector(false, "Getlk")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Lseek(ctx context.Context, off uint64, whence uint32) (uint64, syscall.Errno) {
	fi.Detector(false, "Lseek")
	return 0, syscall.EOPNOTSUPP
}

func (fi FileImplementer) PassthroughFd() (int, bool) {
	fi.Detector(true, "PassthroughFd")
	return 0, false
}

func (fi FileImplementer) Readdirent(ctx context.Context) (*fuse.DirEntry, syscall.Errno) {
	fi.Detector(false, "Readdirent")
	return nil, syscall.EOPNOTSUPP
}

func (fi FileImplementer) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	fi.Detector(false, "Read")
	return nil, syscall.EOPNOTSUPP
}

func (fi FileImplementer) Releasedir(ctx context.Context, releaseFlags uint32) {
	fi.Detector(true, "Releasedir")
}

func (fi FileImplementer) Release(ctx context.Context) syscall.Errno {
	fi.Detector(false, "Release")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Seekdir(ctx context.Context, off uint64) syscall.Errno {
	fi.Detector(false, "Seekdir")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Setattr(ctx context.Context, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	fi.Detector(false, "Setattr")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Setlk(ctx context.Context, owner uint64, lk *fuse.FileLock, flags uint32) syscall.Errno {
	fi.Detector(false, "Setlk")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Setlkw(ctx context.Context, owner uint64, lk *fuse.FileLock, flags uint32) syscall.Errno {
	fi.Detector(false, "Setlkw")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Statx(ctx context.Context, flags uint32, mask uint32, out *fuse.StatxOut) syscall.Errno {
	fi.Detector(false, "Statx")
	return syscall.EOPNOTSUPP
}

func (fi FileImplementer) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	fi.Detector(false, "Write")
	return 0, syscall.EOPNOTSUPP
}
