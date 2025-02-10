package layer

import (
	"context"
	"errors"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

var _ storage.Directory = (*directory)(nil)

type directory struct {
	upperlayer storage.Directory
	underlayer storage.Directory
}

func NewDirectory(upperlayer storage.Directory, underlayer storage.Directory) *directory {
	return &directory{
		upperlayer: upperlayer,
		underlayer: underlayer,
	}
}

// CreateDirectory creates a directory.
func (d *directory) CreateDirectory(ctx context.Context, name string) (storage.Directory, error) {
	var underlayerChild storage.Directory

	// Create the directory on the underlayer
	underlayerChild, err := d.underlayer.CreateDirectory(ctx, name)
	if err != nil {
		return nil, err
	}

	// Create the directory on the upperlayer
	upperlayerChild, err := d.upperlayer.CreateDirectory(ctx, name)
	if err != nil {
		return nil, err
	}

	// Return the new directory
	return NewDirectory(upperlayerChild, underlayerChild), nil
}

// GetInfo returns the directory info.
func (d *directory) GetInfo(_ context.Context) (info.Directory, error) {
	return info.Directory{}, nil
}

// ListFiles returns a map of files.
func (d *directory) ListFiles(ctx context.Context) (map[string]storage.File, error) {
	// Get local files
	upperlayerFiles, err := d.upperlayer.ListFiles(ctx)
	if err != nil {
		return nil, err
	}

	// Get the files
	files := make(map[string]storage.File, len(upperlayerFiles))
	for n := range upperlayerFiles {
		f, err := d.GetFile(ctx, n)
		if err != nil {
			return nil, err
		}
		files[n] = f
	}

	// Get underlayer files
	underlayer, err := d.underlayer.ListFiles(ctx)
	if err != nil {
		return nil, err
	}

	// Merge the two maps
	for n := range underlayer {
		if _, ok := files[n]; ok {
			continue
		}

		f, err := d.GetFile(ctx, n)
		if err != nil {
			return nil, err
		}
		files[n] = f
	}

	return files, nil
}

// GetDirectory returns a child directory.
func (d *directory) GetDirectory(ctx context.Context, name string) (storage.Directory, error) {
	var underlayer storage.Directory
	var err error

	// Get the directory from the underlayer
	underlayer, err = d.underlayer.GetDirectory(ctx, name)
	if err != nil {
		return nil, err
	}

	// Get the directory from the upperlayer
	upperlayer, err := d.upperlayer.GetDirectory(ctx, name)
	if err != nil {
		if !errors.Is(err, storage.ErrDirectoryNotFound) || underlayer == nil {
			return nil, err
		}
	}

	// Return the directory
	return NewDirectory(upperlayer, underlayer), nil
}

// GetFile returns a child file.
func (d *directory) GetFile(ctx context.Context, name string) (storage.File, error) {
	var underlayer storage.File
	var err error

	// Get the directory from the underlayer
	underlayer, err = d.underlayer.GetFile(ctx, name)
	if err != nil {
		return nil, err
	}

	// Get the directory from the upperlayer
	var info info.File
	upperlayerFile, err := d.upperlayer.GetFile(ctx, name)
	if err != nil {
		// If there is an error and it's not a file not found error
		if !errors.Is(err, storage.ErrFileNotFound) {
			return nil, err
		}

		// Get the info from the underlayer
		info, err = underlayer.GetInfo(ctx)
		if err != nil {
			return nil, err
		}

		// Create a new file on the upperlayer
		upperlayerFile, err = d.upperlayer.CreateFile(ctx, name, info)
		if err != nil {
			return nil, err
		}
	} else {
		// Get the info from the upperlayer
		info, err = upperlayerFile.GetInfo(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Return the directory
	return newFile(upperlayerFile, underlayer, info), nil
}

// ListDirectories returns a map of directories.
func (d *directory) ListDirectories(ctx context.Context) (map[string]storage.Directory, error) {
	// Get local directories
	upperlayerDirectories, err := d.upperlayer.ListDirectories(ctx)
	if err != nil {
		return nil, err
	}

	// Get the directories
	directories := make(map[string]storage.Directory, len(upperlayerDirectories))
	for n := range upperlayerDirectories {
		dir, err := d.GetDirectory(ctx, n)
		if err != nil {
			return nil, err
		}
		directories[n] = dir
	}

	// Get underlayer directories
	underlayer, err := d.underlayer.ListDirectories(ctx)
	if err != nil {
		return nil, err
	}

	// Merge the two maps
	for n := range underlayer {
		if _, ok := directories[n]; ok {
			continue
		}

		dir, err := d.GetDirectory(ctx, n)
		if err != nil {
			return nil, err
		}
		directories[n] = dir
	}

	return directories, nil
}

// CreateFile creates a file in the directory.
func (d *directory) CreateFile(ctx context.Context, name string, info info.File) (storage.File, error) {
	var underlayerChild storage.File

	// Create the directory on the underlayer
	underlayerChild, err := d.underlayer.CreateFile(ctx, name, info)
	if err != nil {
		return nil, err
	}

	// Create the file on the upperlayer
	upperlayerFile, err := d.upperlayer.CreateFile(ctx, name, info)
	if err != nil {
		return nil, err
	}

	// Return the new directory
	return newFile(upperlayerFile, underlayerChild, info), nil
}

// RemoveDirectory removes a child directory of the directory.
func (d *directory) RemoveDirectory(ctx context.Context, name string) error {
	// Remove the directory from the underlayer
	if err := d.underlayer.RemoveDirectory(ctx, name); err != nil {
		return err
	}

	// Remove the directory from the upperlayer
	err := d.upperlayer.RemoveDirectory(ctx, name)
	if err == nil || errors.Is(err, storage.ErrDirectoryNotFound) {
		return nil
	}

	return err
}

// RemoveFile removes a child file of the directory.
func (d *directory) RemoveFile(ctx context.Context, name string) error {
	// Remove the directory from the underlayer
	if err := d.underlayer.RemoveFile(ctx, name); err != nil {
		return err
	}

	// Remove the directory from the upperlayer
	err := d.upperlayer.RemoveFile(ctx, name)
	if err == nil || errors.Is(err, storage.ErrFileNotFound) {
		return nil
	}

	return err
}

// RenameFile renames a child file of the directory.
func (d *directory) RenameFile(
	ctx context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// Rename the file on the underlayer
	newParentUnderlayer := newParent.(*directory).underlayer
	if err := d.underlayer.RenameFile(ctx, name, newParentUnderlayer, newName, noReplace); err != nil {
		return err
	}

	// Rename the file on the upperlayer
	newParentBackend := newParent.(*directory).upperlayer
	err := d.upperlayer.RenameFile(ctx, name, newParentBackend, newName, noReplace)
	if err == nil || errors.Is(err, storage.ErrFileNotFound) {
		return nil
	}

	return err
}

// RenameDirectory renames a child directory of the directory.
func (d *directory) RenameDirectory(
	ctx context.Context,
	name string,
	newParent storage.Directory,
	newName string,
	noReplace bool,
) error {
	// Rename the directory on the underlayer
	newParentUnderlayer := newParent.(*directory).underlayer
	if err := d.underlayer.RenameDirectory(ctx, name, newParentUnderlayer, newName, noReplace); err != nil {
		return err
	}

	// Rename the directory on the upperlayer
	newParentBackend := newParent.(*directory).upperlayer
	err := d.upperlayer.RenameDirectory(ctx, name, newParentBackend, newName, noReplace)
	if err == nil || errors.Is(err, storage.ErrDirectoryNotFound) {
		return nil
	}

	return err
}
