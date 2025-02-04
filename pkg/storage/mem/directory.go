package mem

import (
	"context"
	"errors"
	"fmt"

	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.Directory = (*Directory)(nil)

// DirectoryOptions represents the options that can be given to a Directory.
type DirectoryOptions struct {
	Underlayer storage.Directory
}

// Directory is a directory in memory.
type Directory struct {
	directories map[string]*Directory
	files       map[string]*file
	underlayer  storage.Directory
}

// NewDirectory creates a new directory.
func NewDirectory(opts *DirectoryOptions) *Directory {
	d := &Directory{
		directories: make(map[string]*Directory),
		files:       make(map[string]*file),
	}

	if opts != nil {
		d.underlayer = opts.Underlayer
	}

	return d
}

// Underlayer returns the directory underlayer.
func (d *Directory) Underlayer() storage.Directory {
	return d.underlayer
}

// CreateDirectory creates a directory.
func (d *Directory) CreateDirectory(ctx context.Context, name string) (storage.Directory, error) {
	var childUnderlayer storage.Directory
	var err error

	// Check if the directory already exists
	if u := d.Underlayer(); u != nil {
		// If there is an underlayer, then creates it here: it will check if
		// the directory already exists
		childUnderlayer, err = u.CreateDirectory(ctx, name)
		if err != nil {
			return nil, err
		}
	} else if _, exist := d.directories[name]; exist {
		// If already exists, then return an error
		return nil, storage.ErrDirectoryAlreadyExists
	}

	// Create the new directory with its underlayer
	nd := NewDirectory(&DirectoryOptions{
		Underlayer: childUnderlayer,
	})
	d.directories[name] = nd

	return nd, nil
}

// Info returns the directory info.
func (d *Directory) Info(_ context.Context) (storage.DirectoryInfo, error) {
	return storage.DirectoryInfo{}, nil
}

// ListFiles returns a map of files.
func (d *Directory) ListFiles(ctx context.Context) (map[string]storage.File, error) {
	// Get the actual list
	m := make(map[string]storage.File, len(d.directories))
	for p, f := range d.files {
		m[p] = f
	}

	// Check if there is some from under that are not listed here
	if u := d.Underlayer(); u != nil {
		// Get the under dirs
		underFiles, err := u.ListFiles(ctx)
		if err != nil {
			return nil, err
		}

		// Check that they are all present
		for n, uf := range underFiles {
			// If they are, just continue
			if _, ok := m[n]; ok {
				continue
			}

			// Get the info of the underlaying file
			info, err := uf.Info(ctx)
			if err != nil {
				return nil, err
			}

			// Create a new file, based on the underlayer one
			newFile := newFile(info.ChunkSize, &fileOptions{
				Underlayer: uf,
			})

			// Add it to the result and to the memory
			m[n] = newFile
			d.files[n] = newFile
		}
	}

	return m, nil
}

// GetDirectory returns a child directory.
func (d *Directory) GetDirectory(ctx context.Context, name string) (storage.Directory, error) {
	// If it exists, return the directory
	dir, ok := d.directories[name]
	if ok {
		return dir, nil
	}

	// If it doesn't, check that it doesn't exists on underlayer
	if u := d.Underlayer(); u != nil {
		childUnderlayer, err := u.GetDirectory(ctx, name)
		if err != nil {
			return nil, err
		}

		// Add it to the memory if it exists
		dir = NewDirectory(&DirectoryOptions{
			Underlayer: childUnderlayer,
		})
		d.directories[name] = dir

		return dir, nil
	}

	return nil, storage.ErrDirectoryNotExists
}

// GetFile returns a child file.
func (d *Directory) GetFile(ctx context.Context, name string) (storage.File, error) {
	// Check locally
	f, ok := d.files[name]
	if ok {
		return f, nil
	}

	// If it doesn't, check that it doesn't exists on underlayer
	if u := d.Underlayer(); u != nil {
		childUnderlayer, err := u.GetFile(ctx, name)
		if err != nil {
			return nil, err
		}

		// Get info
		fi, err := childUnderlayer.Info(ctx)
		if err != nil {
			return nil, err
		}

		// Add it to the memory if it exists
		file := newFile(fi.ChunkSize, &fileOptions{
			Underlayer:    childUnderlayer,
			ChunkNb:       fi.ChunksCount,
			LastChunkSize: fi.LastChunkSize,
		})
		d.files[name] = file

		return file, nil
	}

	return nil, storage.ErrFileNotExists
}

// ListDirectories returns a map of directories.
func (d *Directory) ListDirectories(ctx context.Context) (map[string]storage.Directory, error) {
	// Get the actual list
	m := make(map[string]storage.Directory, len(d.directories))
	for p, dir := range d.directories {
		m[p] = dir
	}

	// Check if there is some from under that are not listed here
	if u := d.Underlayer(); u != nil {
		// Get the under dirs
		underDirs, err := u.ListDirectories(ctx)
		if err != nil {
			return nil, err
		}

		// Check that they are all present
		for n, ud := range underDirs {
			// If they are, just continue
			if _, ok := m[n]; ok {
				continue
			}

			// Create a new directory, based on the underlayer one
			newDir := NewDirectory(&DirectoryOptions{
				Underlayer: ud,
			})

			// Add it to the result and to the memory
			m[n] = newDir
			d.directories[n] = newDir
		}
	}

	return m, nil
}

// CreateFile creates a file in the directory.
func (d *Directory) CreateFile(ctx context.Context, name string, chunkSize int) (storage.File, error) {
	var childUnderlayer storage.File
	var err error

	// Check if the file already exists
	if u := d.Underlayer(); u != nil {
		// Create the file in the underlayer, that will check that the file already exists
		childUnderlayer, err = u.CreateFile(ctx, name, chunkSize)
		if err != nil {
			return nil, err
		}
	} else if _, exist := d.files[name]; exist {
		// If the file already exists, return an error
		return nil, fmt.Errorf("couldn't create file %q: %w", name, storage.ErrFileAlreadyExists)
	}

	f := newFile(chunkSize, &fileOptions{
		Underlayer: childUnderlayer,
	})
	d.files[name] = f
	return f, nil
}

// RemoveDirectory removes a child directory of the directory.
func (d *Directory) RemoveDirectory(ctx context.Context, name string) error {
	// Remove in underlayer
	if u := d.Underlayer(); u != nil {
		if err := u.RemoveDirectory(ctx, name); err != nil {
			return err
		}
	}

	// Actually delete the directory
	delete(d.directories, name)

	return nil
}

// RemoveFile removes a child file of the directory.
func (d *Directory) RemoveFile(ctx context.Context, name string) error {
	// Check if the file exists
	if u := d.Underlayer(); u != nil {
		// If there is an underlayer, remove it there, it will check if the file exists
		if err := u.RemoveFile(ctx, name); err != nil {
			return err
		}
	} else if _, exist := d.files[name]; !exist {
		// If the file doesn't exists, return an error
		return storage.ErrFileNotExists
	}
	delete(d.files, name)
	return nil
}

func (d *Directory) checkIfFileOrDirectoryAlreadyExists(name string) error {
	if _, exist := d.directories[name]; exist {
		return fmt.Errorf("%w: %q", storage.ErrDirectoryAlreadyExists, name)
	}
	if _, exist := d.files[name]; exist {
		return fmt.Errorf("%w: %q", storage.ErrFileAlreadyExists, name)
	}
	return nil
}

// RenameFile renames a child file of the directory.
func (d *Directory) RenameFile(
	ctx context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// Get the the file
	f, err := d.GetFile(ctx, name)
	if err != nil {
		return err
	}

	// Check if it doesn't not exist already
	if err := newParent.(*Directory).checkIfFileOrDirectoryAlreadyExists(newName); err != nil {
		if errors.Is(err, storage.ErrFileAlreadyExists) && noReplace {
			// If noReplace is set and the file already exists, return error
			return err
		} else if !errors.Is(err, storage.ErrFileAlreadyExists) {
			// If another error, then return
			return err
		}
	}

	// If there is an underlayer, then rename the file here first
	if u := d.Underlayer(); u != nil {
		if err := u.RenameFile(ctx, name, newParent.Underlayer(), newName, noReplace); err != nil {
			return err
		}
	}

	// Add it to new parent and remove it from current parent
	newParent.(*Directory).files[newName] = f.(*file)
	delete(d.files, name)

	return nil
}

// RenameDirectory renames a child directory of the directory.
func (d *Directory) RenameDirectory(
	ctx context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// Get directory
	dir, err := d.GetDirectory(ctx, name)
	if err != nil {
		return err
	}

	// Check if it doesn't not exist already
	if err := newParent.(*Directory).checkIfFileOrDirectoryAlreadyExists(newName); err != nil {
		if errors.Is(err, storage.ErrDirectoryAlreadyExists) && noReplace {
			// If noReplace is set and the file already exists, return error
			return err
		} else if !errors.Is(err, storage.ErrDirectoryAlreadyExists) {
			// If another error, then return
			return err
		}
	}

	// If there is an underlayer, then rename the directory here first
	if u := d.Underlayer(); u != nil {
		if err := u.RenameDirectory(ctx, name, newParent.Underlayer(), newName, noReplace); err != nil {
			return err
		}
	}

	// Add it to new parent and remove it from current parent
	newParent.(*Directory).directories[newName] = dir.(*Directory)
	delete(d.directories, name)

	return nil
}
