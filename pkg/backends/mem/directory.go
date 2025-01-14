package mem

import (
	"context"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/backends"
)

type directory struct {
	attr  fuse.Attr
	dirs  map[string]*directory
	files map[string]*file
}

func newEmptyDirectory() *directory {
	return &directory{
		dirs:  make(map[string]*directory),
		files: make(map[string]*file),
	}
}

func (dir *directory) checkIfFileOrDirectoryAlreadyExists(name string) syscall.Errno {
	// Check in directories
	if _, ok := dir.dirs[name]; ok {
		return syscall.EEXIST
	}

	// Check in files
	if _, ok := dir.files[name]; ok {
		return syscall.EEXIST
	}

	return fs.OK
}

func (dir *directory) CreateDirectory(ctx context.Context, name string) (backends.Directory, syscall.Errno) {
	// Check if it doesn't not exist already
	if errno := dir.checkIfFileOrDirectoryAlreadyExists(name); errno != fs.OK {
		return nil, errno
	}

	// Create a new directory
	c := newEmptyDirectory()

	// Add it to childs
	dir.dirs[name] = c

	return c, fs.OK
}

func (dir *directory) GetDirectory(ctx context.Context, name string) (backends.Directory, syscall.Errno) {
	// Check if this is not already a file
	if _, ok := dir.files[name]; ok {
		return nil, syscall.ENOTDIR
	}

	// Get and check if it exists
	d, ok := dir.dirs[name]
	if !ok {
		return nil, syscall.ENOENT
	}

	return d, fs.OK
}

func (dir *directory) GetFile(ctx context.Context, name string) (backends.File, syscall.Errno) {
	// Get and check if it exists
	f, ok := dir.files[name]
	if !ok {
		return nil, syscall.ENOENT
	}

	return f, fs.OK
}

func (dir *directory) ListEntries(ctx context.Context) ([]fuse.DirEntry, syscall.Errno) {
	list := make([]fuse.DirEntry, 0, len(dir.dirs)+len(dir.files))

	// Add directories
	for n := range dir.dirs {
		list = append(list, fuse.DirEntry{
			Mode: fuse.S_IFDIR,
			Name: n,
			// TODO: add Ino
		})
	}

	// Add files
	for n := range dir.files {
		list = append(list, fuse.DirEntry{
			Mode: fuse.S_IFREG,
			Name: n,
			// TODO: add Ino
		})
	}

	return list, fs.OK
}

func (dir *directory) CreateFile(ctx context.Context, name string) (backends.File, syscall.Errno) {
	// Check if it doesn't not exist already
	if errno := dir.checkIfFileOrDirectoryAlreadyExists(name); errno != fs.OK {
		return nil, errno
	}

	// Create file
	f := newEmptyFile()

	// Add it to children
	dir.files[name] = f

	return f, fs.OK
}

func (dir *directory) RemoveDirectory(ctx context.Context, name string) syscall.Errno {
	// Check if it exists
	if _, ok := dir.dirs[name]; !ok {
		return syscall.ENOENT
	}

	// Remove it from memory
	delete(dir.dirs, name)

	return fs.OK
}

func (dir *directory) RemoveFile(ctx context.Context, name string) syscall.Errno {
	// Check if it exists
	if _, ok := dir.files[name]; !ok {
		return syscall.ENOENT
	}

	// Remove it from memory
	delete(dir.files, name)

	return fs.OK
}

func (dir *directory) GetAttributes(ctx context.Context, attr *fuse.Attr) syscall.Errno {
	*attr = dir.attr
	return fs.OK
}

func (dir *directory) SetAttributes(ctx context.Context, in *fuse.SetAttrIn) syscall.Errno {
	// TODO
	return fs.OK
}
