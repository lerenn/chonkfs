package mem

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.Directory = (*Directory)(nil)

type Directory struct {
	directories map[string]*Directory
	files       map[string]*File
}

func NewDirectory() *Directory {
	return &Directory{
		directories: make(map[string]*Directory),
		files:       make(map[string]*File),
	}
}

func (d *Directory) CreateDirectory(ctx context.Context, name string) (storage.Directory, error) {
	if _, exist := d.directories[name]; exist {
		return nil, storage.ErrDirectoryAlreadyExists
	}

	nd := NewDirectory()
	d.directories[name] = nd
	return nd, nil
}

func (d *Directory) Info(ctx context.Context) (storage.DirectoryInfo, error) {
	return storage.DirectoryInfo{}, nil
}

func (d *Directory) ListFiles(ctx context.Context) (map[string]storage.File, error) {
	m := make(map[string]storage.File, len(d.files))
	for p, f := range d.files {
		m[p] = f
	}
	return m, nil
}

func (d *Directory) GetDirectory(ctx context.Context, name string) (storage.Directory, error) {
	dir, ok := d.directories[name]
	if !ok {
		return nil, storage.ErrDirectoryNotExists
	}
	return dir, nil
}

func (d *Directory) GetFile(ctx context.Context, name string) (storage.File, error) {
	f, ok := d.files[name]
	if !ok {
		return nil, storage.ErrFileNotExists
	}
	return f, nil
}

func (d *Directory) ListDirectories(ctx context.Context) (map[string]storage.Directory, error) {
	m := make(map[string]storage.Directory, len(d.directories))
	for p, dir := range d.directories {
		m[p] = dir
	}
	return m, nil
}

func (d *Directory) CreateFile(ctx context.Context, name string, chunkSize int) (storage.File, error) {
	if _, exist := d.files[name]; exist {
		return nil, storage.ErrFileAlreadyExists
	}

	f := newFile(chunkSize)
	d.files[name] = f
	return f, nil
}

func (d *Directory) RemoveDirectory(ctx context.Context, name string) error {
	if _, exist := d.directories[name]; !exist {
		return storage.ErrDirectoryNotExists
	}
	delete(d.directories, name)
	return nil
}

func (d *Directory) RemoveFile(ctx context.Context, name string) error {
	if _, exist := d.files[name]; !exist {
		return storage.ErrFileNotExists
	}
	delete(d.files, name)
	return nil
}

func (d *Directory) checkIfFileOrDirectoryAlreadyExists(name string) error {
	if _, exist := d.directories[name]; exist {
		return storage.ErrDirectoryAlreadyExists
	}
	if _, exist := d.files[name]; exist {
		return storage.ErrFileAlreadyExists
	}
	return nil
}

func (d *Directory) RenameFile(ctx context.Context, name string, newParent storage.Directory, newName string) error {
	// Get the directory or the file
	f, fileExist := d.files[name]
	if !fileExist {
		return storage.ErrFileNotExists
	}

	// Check if it doesn't not exist already
	if err := newParent.(*Directory).checkIfFileOrDirectoryAlreadyExists(newName); err != nil {
		return err
	}

	// Add it to new parent and remove it from current parent
	newParent.(*Directory).files[newName] = f
	delete(d.files, name)

	return nil
}

func (d *Directory) RenameDirectory(ctx context.Context, name string, newParent storage.Directory, newName string) error {
	// Get the directory or the file
	dir, dirExist := d.directories[name]
	if !dirExist {
		return storage.ErrFileNotExists
	}

	// Check if it doesn't not exist already
	if err := newParent.(*Directory).checkIfFileOrDirectoryAlreadyExists(newName); err != nil {
		return err
	}

	// Add it to new parent and remove it from current parent
	newParent.(*Directory).directories[newName] = dir
	delete(d.directories, name)

	return nil
}
