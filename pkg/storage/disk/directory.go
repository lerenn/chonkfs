package disk

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

type directory struct {
	path string
}

// NewDirectory creates a new directory representation.
func NewDirectory(path string) storage.Directory {
	return &directory{
		path: path,
	}
}

func (d *directory) getChildPath(name string) string {
	return fmt.Sprintf("%s/%s", d.path, name)
}

func (d *directory) getChildMetadataPath(name string) string {
	return fmt.Sprintf("%s/%s", d.getChildPath(name), metadataFileName)
}

func (d *directory) ensureChildDoesNotExists(name string) error {
	path := d.getChildPath(name)
	metadataPath := d.getChildMetadataPath(name)

	_, err := os.Stat(path)
	if err == nil {
		// Check if there is a metadata file
		_, err = os.Stat(metadataPath)
		if err == nil {
			return fmt.Errorf("%w: %q", storage.ErrFileAlreadyExists, name)
		}

		return fmt.Errorf("%w: %q", storage.ErrDirectoryAlreadyExists, name)
	} else if !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (d *directory) writeChildMetadata(name string, info info.File) error {
	return writeMetadata(path.Join(d.path, name), info)
}

func (d *directory) readChildMetadata(name string) (info.File, error) {
	return readMetadata(path.Join(d.path, name))
}

// CreateDirectory creates a directory.
func (d *directory) CreateDirectory(_ context.Context, name string) (storage.Directory, error) {
	// Check if a file or a directory exists
	if err := d.ensureChildDoesNotExists(name); err != nil {
		return nil, err
	}

	// Create directory and return representation
	path := d.getChildPath(name)
	return NewDirectory(path), os.Mkdir(path, 0755)
}

// GetDirectory returns a child directory of the directory.
func (d *directory) GetDirectory(_ context.Context, name string) (storage.Directory, error) {
	path := d.getChildPath(name)

	// Check if directory exists
	err := d.ensureChildDoesNotExists(name)
	switch {
	case err == nil:
		return nil, fmt.Errorf("%w: %q", storage.ErrDirectoryNotFound, name)
	case errors.Is(err, storage.ErrFileAlreadyExists):
		return nil, fmt.Errorf("%w: %q", storage.ErrIsFile, name)
	case !errors.Is(err, storage.ErrDirectoryAlreadyExists):
		return nil, err
	}

	// Return representation
	return NewDirectory(path), nil
}

// GetInfo returns the directory information.
func (d *directory) GetInfo(_ context.Context) (info.Directory, error) {
	return info.Directory{}, nil
}

// CreateFile creates a file.
func (d *directory) CreateFile(_ context.Context, name string, info info.File) (storage.File, error) {
	path := d.getChildPath(name)

	// Check if there is a file with this name
	if err := d.ensureChildDoesNotExists(name); err != nil {
		return nil, err
	}

	// Create file representation
	f, err := newFile(path, info)
	if err != nil {
		return nil, err
	}

	// Create directory representing the file
	if err := os.Mkdir(path, 0755); err != nil {
		return nil, err
	}

	return f, d.writeChildMetadata(name, info)
}

// GetFile returns a child file.
func (d *directory) GetFile(_ context.Context, name string) (storage.File, error) {
	// Check if file exists
	err := d.ensureChildDoesNotExists(name)
	switch {
	case err == nil:
		return nil, fmt.Errorf("%w: %q", storage.ErrFileNotFound, name)
	case errors.Is(err, storage.ErrDirectoryAlreadyExists):
		return nil, fmt.Errorf("%w: %q", storage.ErrIsDirectory, name)
	case !errors.Is(err, storage.ErrFileAlreadyExists):
		return nil, err
	}

	// Read metadata
	info, err := d.readChildMetadata(name)
	if err != nil {
		return nil, err
	}

	// Create file representation
	path := d.getChildPath(name)
	return newFile(path, info)
}

// ListFiles returns a map of files.
func (d *directory) ListFiles(_ context.Context) (map[string]storage.File, error) {
	entries, err := os.ReadDir(d.path)
	if err != nil {
		return nil, err
	}

	files := make(map[string]storage.File)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Read metadata
		info, err := d.readChildMetadata(entry.Name())
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		} else if os.IsNotExist(err) {
			// This is a directory
			continue
		}

		// Create file representation
		path := path.Join(d.path, entry.Name())
		f, err := newFile(path, info)
		if err != nil {
			return nil, err
		}

		files[entry.Name()] = f
	}

	return files, nil
}

// RemoveDirectory removes a directory.
func (d *directory) RemoveDirectory(_ context.Context, name string) error {
	// Check if directory exists
	err := d.ensureChildDoesNotExists(name)
	switch {
	case err == nil:
		return fmt.Errorf("%w: %q", storage.ErrDirectoryNotFound, name)
	case errors.Is(err, storage.ErrFileAlreadyExists):
		return fmt.Errorf("%w: %q", storage.ErrIsFile, name)
	case !errors.Is(err, storage.ErrDirectoryAlreadyExists):
		return err
	}

	// Remove directory
	return os.Remove(d.getChildPath(name))
}

// ListDirectories returns a map of directories.
func (d *directory) ListDirectories(_ context.Context) (map[string]storage.Directory, error) {
	entries, err := os.ReadDir(d.path)
	if err != nil {
		return nil, err
	}

	directories := make(map[string]storage.Directory)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check if there is metadata
		metadataPath := d.getChildMetadataPath(entry.Name())
		_, err = os.Stat(metadataPath)
		if err == nil {
			// This is a file
			continue
		} else if !os.IsNotExist(err) {
			return nil, err
		}

		// Create directory representation
		path := path.Join(d.path, entry.Name())
		directories[entry.Name()] = NewDirectory(path)
	}

	return directories, nil
}

// RemoveFile removes a file.
func (d *directory) RemoveFile(_ context.Context, name string) error {
	// Check if file or directory exists
	err := d.ensureChildDoesNotExists(name)
	switch {
	case err == nil:
		return fmt.Errorf("%w: %q", storage.ErrFileNotFound, name)
	case errors.Is(err, storage.ErrDirectoryAlreadyExists):
		return fmt.Errorf("%w: %q", storage.ErrIsDirectory, name)
	case !errors.Is(err, storage.ErrFileAlreadyExists):
		return err
	}

	return os.RemoveAll(d.getChildPath(name))
}

// RenameFile renames a file.
func (d *directory) RenameFile(
	_ context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// Check if file exists
	err := d.ensureChildDoesNotExists(name)
	switch {
	case err == nil:
		return fmt.Errorf("%w: %q", storage.ErrFileNotFound, name)
	case errors.Is(err, storage.ErrDirectoryAlreadyExists):
		return fmt.Errorf("%w: %q", storage.ErrIsDirectory, name)
	case !errors.Is(err, storage.ErrFileAlreadyExists):
		return err
	}

	// Check if there is a file with the new name
	newPath := newParent.(*directory).getChildPath(newName)
	err = newParent.(*directory).ensureChildDoesNotExists(newName)
	switch {
	case err == nil:
		// Nothing to do
	case !noReplace && (errors.Is(err, storage.ErrDirectoryAlreadyExists) ||
		errors.Is(err, storage.ErrFileAlreadyExists)):
		if err := os.RemoveAll(newPath); err != nil {
			return err
		}
	default:
		return err
	}

	// Move the file
	oldPath := d.getChildPath(name)
	return os.Rename(oldPath, newPath)
}

// RenameDirectory renames a directory.
func (d *directory) RenameDirectory(
	_ context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// Check if directory exists
	err := d.ensureChildDoesNotExists(name)
	switch {
	case err == nil:
		return fmt.Errorf("%w: %q", storage.ErrDirectoryNotFound, name)
	case errors.Is(err, storage.ErrFileAlreadyExists):
		return fmt.Errorf("%w: %q", storage.ErrIsFile, name)
	case !errors.Is(err, storage.ErrDirectoryAlreadyExists):
		return err
	}

	// Check if there is a directory with the new name
	newPath := newParent.(*directory).getChildPath(newName)
	err = newParent.(*directory).ensureChildDoesNotExists(newName)
	switch {
	case err == nil:
		// Nothing to do
	case !noReplace && (errors.Is(err, storage.ErrDirectoryAlreadyExists) ||
		errors.Is(err, storage.ErrFileAlreadyExists)):
		if err := os.RemoveAll(newPath); err != nil {
			return err
		}
	default:
		return err
	}

	// Move the directory
	oldPath := d.getChildPath(name)
	return os.Rename(oldPath, newPath)
}
