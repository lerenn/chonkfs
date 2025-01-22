package chonker

import (
	"context"
	"io"
	"log"

	"github.com/hanwen/go-fuse/v2/fuse"
)

type DirectoryOption func(dir *directory)

func WithDirectoryLogger(logger *log.Logger) DirectoryOption {
	return func(dir *directory) {
		dir.logger = logger
	}
}

var _ Directory = (*directory)(nil)

type directory struct {
	dirs   map[string]*directory
	files  map[string]*file
	logger *log.Logger
	opts   []DirectoryOption
}

func NewDirectory(opts ...DirectoryOption) *directory {
	// Create a default directory
	d := &directory{
		dirs:   make(map[string]*directory),
		files:  make(map[string]*file),
		logger: log.New(io.Discard, "", 0),
		opts:   opts,
	}

	// Apply options
	for _, opt := range opts {
		opt(d)
	}

	return d
}

func (dir *directory) checkIfFileOrDirectoryAlreadyExists(name string) error {
	// Check in directories
	if _, ok := dir.dirs[name]; ok {
		return ErrAlreadyExists
	}

	// Check in files
	if _, ok := dir.files[name]; ok {
		return ErrAlreadyExists
	}

	return nil
}

func (dir *directory) CreateDirectory(ctx context.Context, name string) (Directory, error) {
	// Check if it doesn't not exist already
	if err := dir.checkIfFileOrDirectoryAlreadyExists(name); err != nil {
		return nil, err
	}

	// Create a new directory
	c := NewDirectory(dir.opts...)

	// Add it to childs
	dir.dirs[name] = c

	return c, nil
}

func (dir *directory) GetDirectory(ctx context.Context, name string) (Directory, error) {
	// Check if this is not already a file
	if _, ok := dir.files[name]; ok {
		return nil, ErrNotDirectory
	}

	// Get and check if it exists
	d, ok := dir.dirs[name]
	if !ok {
		return nil, ErrNoEntry
	}

	return d, nil
}

func (dir *directory) GetFile(ctx context.Context, name string) (File, error) {
	// Get and check if it exists
	f, ok := dir.files[name]
	if !ok {
		return nil, ErrNoEntry
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

func (dir *directory) CreateFile(ctx context.Context, name string, chunkSize int) (File, error) {
	// Check if it doesn't not exist already
	if err := dir.checkIfFileOrDirectoryAlreadyExists(name); err != nil {
		return nil, err
	}

	// Create file
	f := newFile(chunkSize,
		WithFileLogger(dir.logger))

	// Add it to children
	dir.files[name] = f

	return f, nil
}

func (dir *directory) RemoveDirectory(ctx context.Context, name string) error {
	// Check if it exists
	if _, ok := dir.dirs[name]; !ok {
		return ErrNoEntry
	}

	// Remove it from memory
	delete(dir.dirs, name)

	return nil
}

func (dir *directory) RemoveFile(ctx context.Context, name string) error {
	// Check if it exists
	if _, ok := dir.files[name]; !ok {
		return ErrNoEntry
	}

	// Remove it from memory
	delete(dir.files, name)

	return nil
}

func (dir *directory) GetAttributes(ctx context.Context) (fuse.Attr, error) {
	// TODO
	return fuse.Attr{}, nil
}

func (dir *directory) SetAttributes(ctx context.Context, in *fuse.SetAttrIn) error {
	// TODO
	return nil
}

func (dir *directory) RenameEntry(ctx context.Context, name string, newParent Directory, newName string) error {
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
		return ErrNoEntry
	}

	return nil
}
