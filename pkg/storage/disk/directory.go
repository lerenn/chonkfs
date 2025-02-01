package disk

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/lerenn/chonkfs/pkg/storage"
)

const (
	fileInfoName   = ".file"
	defaultDirMode = os.FileMode(01755)
)

var _ storage.Directory = (*Directory)(nil)

// DirectoryOptions represents the options that can be given to a Directory.
type DirectoryOptions struct {
	Underlayer storage.Directory
}

// Directory is a directory on disk.
type Directory struct {
	path string
	opts *DirectoryOptions
}

// NewDirectory creates a new directory.
func NewDirectory(path string, opts *DirectoryOptions) *Directory {
	return &Directory{
		path: path,
		opts: opts,
	}
}

// Underlayer returns the directory underlayer.
func (d *Directory) Underlayer() storage.Directory {
	if d.opts == nil {
		return nil
	}

	return d.opts.Underlayer
}

func (d *Directory) createLocalDirectory(name string, underlayer storage.Directory) (storage.Directory, error) {
	// Create a directory on disk
	childPath := path.Join(d.path, name)
	if err := os.Mkdir(childPath, defaultDirMode); err != nil {
		return nil, err
	}

	// Create a new directory
	return NewDirectory(childPath, &DirectoryOptions{
		Underlayer: underlayer,
	}), nil
}

// CreateDirectory creates a directory.
func (d *Directory) CreateDirectory(ctx context.Context, name string) (storage.Directory, error) {
	var childUnderlayer storage.Directory
	var err error

	// Create the directory on the underlayer first
	if u := d.Underlayer(); u != nil {
		// If there is an underlayer, then creates it here: it will check if
		// the directory already exists
		childUnderlayer, err = u.CreateDirectory(ctx, name)
		if err != nil {
			return nil, err
		}
	}

	// Create local directory
	return d.createLocalDirectory(name, childUnderlayer)
}

// Info returns the directory info.
func (d *Directory) Info(_ context.Context) (storage.DirectoryInfo, error) {
	return storage.DirectoryInfo{}, nil
}

func (d *Directory) listLocalFiles(ctx context.Context) (map[string]storage.File, error) {
	// Get the entries
	entries, err := os.ReadDir(d.path)
	if err != nil {
		return nil, err
	}

	// Get the entries that are files
	files := make(map[string]storage.File)
	for _, e := range entries {
		name := e.Name()

		// Check if this is not a dir from system perspective, continue if it is not
		if !e.IsDir() {
			continue
		}

		// Check if there is a file info in the directory
		fileInfoPath := path.Join(d.path, name, fileInfoName)
		_, err := os.Stat(fileInfoPath)
		if os.IsNotExist(err) {
			// If it doesn't exists, it means that this is a true directory from
			// chonkfs perspective
			continue
		} else if err != nil {
			// If this is another error, other than doesn't exist, stop here
			return nil, err
		}

		// Get corresponding directory
		file, err := d.GetFile(ctx, name)
		if err != nil {
			return nil, err
		}

		files[name] = file
	}

	return files, nil
}

// ListFiles returns a map of files.
func (d *Directory) ListFiles(ctx context.Context) (map[string]storage.File, error) {
	// Get local directories
	files, err := d.listLocalFiles(ctx)
	if err != nil {
		return nil, err
	}

	// Get underlayer directories
	if u := d.Underlayer(); u != nil {
		underFiles, err := u.ListFiles(ctx)
		if err != nil {
			return nil, err
		}

		// Check that all directories are present
		for n, uf := range underFiles {
			// Continue if they are present
			if _, ok := files[n]; ok {
				continue
			}

			// Get info from file
			fi, err := uf.Info(ctx)
			if err != nil {
				return nil, err
			}

			// Create it if they are not
			nf, err := d.createLocalFile(n, fi.ChunkSize, uf)
			if err != nil {
				return nil, err
			}

			files[n] = nf
		}
	}

	return files, nil
}

// GetDirectory returns a child directory.
func (d *Directory) GetDirectory(ctx context.Context, name string) (storage.Directory, error) {
	var childUnderlayer storage.Directory
	var err error

	// Check the directory exists
	childPath := path.Join(d.path, name)
	if u := d.Underlayer(); u != nil {
		// If there is an underlayer, then gets it here, it will check for it
		childUnderlayer, err = u.GetDirectory(ctx, name)
		if err != nil {
			return nil, err
		}
	} else if _, err := os.Stat(childPath); os.IsNotExist(err) {
		// If there is no underlayer and no directory, then return error
		return nil, storage.ErrDirectoryNotExists
	} else if err != nil {
		// If there is no underlayer and an error, return it
		return nil, err
	}

	return NewDirectory(childPath, &DirectoryOptions{
		Underlayer: childUnderlayer,
	}), nil
}

func (d *Directory) getFileLocally(path string, underlayer storage.File) (storage.File, error) {
	// Get info
	fi, err := readFileInfo(path)
	if err != nil {
		return nil, err
	}

	return newFile(path, fi.ChunkSize, &fileOptions{
		Underlayer: underlayer,
	}), nil
}

// GetFile returns a child file.
func (d *Directory) GetFile(ctx context.Context, name string) (storage.File, error) {
	var childUnderlayer storage.File
	var err error

	// Check the directory exists
	childPath := path.Join(d.path, name)
	if u := d.Underlayer(); u != nil {
		// If there is an underlayer, then gets it here, it will check for it
		childUnderlayer, err = u.GetFile(ctx, name)
		if err != nil {
			return nil, err
		}
	}

	// Check the file exists
	_, err = os.Stat(childPath)
	if err == nil {
		return d.getFileLocally(childPath, childUnderlayer)
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// Get info from underlayer
	info, err := childUnderlayer.Info(ctx)
	if err != nil {
		return nil, err
	}

	return newFile(childPath, info.ChunkSize, &fileOptions{
		Underlayer: childUnderlayer,
	}), nil
}

func (d *Directory) listLocalDirectories(ctx context.Context) (map[string]storage.Directory, error) {
	// Get the entries
	entries, err := os.ReadDir(d.path)
	if err != nil {
		return nil, err
	}

	// Get the entries that are directories
	dirs := make(map[string]storage.Directory)
	for _, e := range entries {
		name := e.Name()

		// Check if this is a dir, otherwise continue
		if !e.IsDir() {
			continue
		}

		// Check if there is a file info in the directory
		fileInfoPath := path.Join(d.path, name, fileInfoName)
		_, err := os.Stat(fileInfoPath)
		if err == nil {
			// Skip as this is a file from chonkfs perspective
			continue
		} else if !os.IsNotExist(err) {
			// If this is another error, other than doesn't exist, stop here
			return nil, err
		}

		// Get corresponding directory
		dir, err := d.GetDirectory(ctx, name)
		if err != nil {
			return nil, err
		}

		dirs[name] = dir
	}

	return dirs, nil
}

// ListDirectories returns a map of directories.
func (d *Directory) ListDirectories(ctx context.Context) (map[string]storage.Directory, error) {
	// Get local directories
	dirs, err := d.listLocalDirectories(ctx)
	if err != nil {
		return nil, err
	}

	// Get underlayer directories
	if u := d.Underlayer(); u != nil {
		underDirs, err := u.ListDirectories(ctx)
		if err != nil {
			return nil, err
		}

		// Check that all directories are present
		for n, ud := range underDirs {
			// Continue if they are present
			if _, ok := dirs[n]; ok {
				continue
			}

			// Create it if they are not
			nd, err := d.createLocalDirectory(n, ud)
			if err != nil {
				return nil, err
			}

			dirs[n] = nd
		}
	}

	return dirs, nil
}

func (d *Directory) createLocalFile(
	name string,
	chunkSize int,
	underlayer storage.File,
) (storage.File, error) {
	// Create a directory on disk that represent the file
	childPath := path.Join(d.path, name)
	if err := os.Mkdir(childPath, defaultDirMode); err != nil {
		return nil, err
	}

	// Create a new file
	f := newFile(childPath, chunkSize, &fileOptions{
		Underlayer: underlayer,
	})

	// Write the info and return
	return f, writeFileInfo(storage.FileInfo{
		ChunkSize: chunkSize,
	}, childPath)
}

// CreateFile creates a file in the directory.
func (d *Directory) CreateFile(ctx context.Context, name string, chunkSize int) (storage.File, error) {
	var childUnderlayer storage.File
	var err error

	// Create file on the underlayer if it exists
	if u := d.Underlayer(); u != nil {
		childUnderlayer, err = u.CreateFile(ctx, name, chunkSize)
		if err != nil {
			return nil, err
		}
	}

	// Create the file on the disk
	return d.createLocalFile(name, chunkSize, childUnderlayer)
}

// RemoveDirectory removes a child directory of the directory.
func (d *Directory) RemoveDirectory(ctx context.Context, name string) error {
	// Remove in underlayer
	if u := d.Underlayer(); u != nil {
		if err := u.RemoveDirectory(ctx, name); err != nil {
			return err
		}
	}

	// Remove on disk
	return os.RemoveAll(path.Join(d.path, name))
}

// RemoveFile removes a child file of the directory.
func (d *Directory) RemoveFile(ctx context.Context, name string) error {
	// Remove in underlayer
	if u := d.Underlayer(); u != nil {
		if err := u.RemoveFile(ctx, name); err != nil {
			return err
		}
	}

	// Remove on disk
	return os.RemoveAll(path.Join(d.path, name))
}

func (d *Directory) checkIfFileOrDirectoryAlreadyExists(_ string) error {
	return fmt.Errorf("not implemented")
}

// RenameFile renames a child file of the directory.
func (d *Directory) RenameFile(
	_ context.Context,
	_ string,
	_ storage.Directory,
	_ string,
	_ bool,
) error {
	return fmt.Errorf("not implemented")
}

// RenameDirectory renames a child directory of the directory.
func (d *Directory) RenameDirectory(
	ctx context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// If there is an underlayer, then rename the directory here first
	if u := d.Underlayer(); u != nil {
		if err := u.RenameDirectory(ctx, name, newParent.Underlayer(), newName, noReplace); err != nil {
			return err
		}
	}

	// Rename on the disk
	oldPath := path.Join(d.path, name)
	newPath := path.Join(newParent.(*Directory).path, newName)
	return os.Rename(oldPath, newPath)
}
