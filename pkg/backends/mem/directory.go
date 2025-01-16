package mem

import (
	"context"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/lerenn/chonkfs/pkg/backends"
)

var _ backends.Directory = (*directory)(nil)

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

func (dir *directory) checkIfFileOrDirectoryAlreadyExists(name string) error {
	// Check in directories
	if _, ok := dir.dirs[name]; ok {
		return backends.ErrAlreadyExists
	}

	// Check in files
	if _, ok := dir.files[name]; ok {
		return backends.ErrAlreadyExists
	}

	return nil
}

func (dir *directory) CreateDirectory(ctx context.Context, name string) (backends.Directory, error) {
	// Check if it doesn't not exist already
	if err := dir.checkIfFileOrDirectoryAlreadyExists(name); err != nil {
		return nil, err
	}

	// Create a new directory
	c := newEmptyDirectory()

	// Add it to childs
	dir.dirs[name] = c

	return c, nil
}

func (dir *directory) GetDirectory(ctx context.Context, name string) (backends.Directory, error) {
	// Check if this is not already a file
	if _, ok := dir.files[name]; ok {
		return nil, backends.ErrNotDirectory
	}

	// Get and check if it exists
	d, ok := dir.dirs[name]
	if !ok {
		return nil, backends.ErrNoEntry
	}

	return d, nil
}

func (dir *directory) GetFile(ctx context.Context, name string) (backends.File, error) {
	// Get and check if it exists
	f, ok := dir.files[name]
	if !ok {
		return nil, backends.ErrNoEntry
	}

	return f, nil
}

func (dir *directory) ListEntries(ctx context.Context) ([]fuse.DirEntry, error) {
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

	return list, nil
}

func (dir *directory) CreateFile(ctx context.Context, name string) (backends.File, error) {
	// Check if it doesn't not exist already
	if err := dir.checkIfFileOrDirectoryAlreadyExists(name); err != nil {
		return nil, err
	}

	// Create file
	f := newEmptyFile()

	// Add it to children
	dir.files[name] = f

	return f, nil
}

func (dir *directory) RemoveDirectory(ctx context.Context, name string) error {
	// Check if it exists
	if _, ok := dir.dirs[name]; !ok {
		return backends.ErrNoEntry
	}

	// Remove it from memory
	delete(dir.dirs, name)

	return nil
}

func (dir *directory) RemoveFile(ctx context.Context, name string) error {
	// Check if it exists
	if _, ok := dir.files[name]; !ok {
		return backends.ErrNoEntry
	}

	// Remove it from memory
	delete(dir.files, name)

	return nil
}

func (dir *directory) GetAttributes(ctx context.Context) (fuse.Attr, error) {
	return dir.attr, nil
}

func (dir *directory) SetAttributes(ctx context.Context, in *fuse.SetAttrIn) error {
	// TODO
	return nil
}

func (dir *directory) RenameNode(ctx context.Context, name string, newParent backends.Directory, newName string) error {
	// Get the directory or the file
	d, dirExist := dir.dirs[name]
	f, fileExist := dir.files[name]

	// Check if it doesn't not exist already
	if err := newParent.(*directory).checkIfFileOrDirectoryAlreadyExists(newName); err != nil {
		return err
	}

	// Add it to new parent and remove it from current parent
	switch {
	case dirExist:
		newParent.(*directory).dirs[newName] = d
		delete(dir.dirs, name)
	case fileExist:
		newParent.(*directory).files[newName] = f
		delete(dir.files, name)
	default:
		return backends.ErrNoEntry
	}

	return nil
}
