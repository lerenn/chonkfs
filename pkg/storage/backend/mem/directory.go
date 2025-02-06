package mem

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/lerenn/chonkfs/pkg/storage/backend"
)

type directory struct {
	directories map[string]*directory
	files       map[string]*file
}

func newDirectory() *directory {
	return &directory{
		directories: make(map[string]*directory),
		files:       make(map[string]*file),
	}
}

func (d *directory) CreateDirectory(path string) error {
	// Split path
	parts := strings.Split(path, "/")

	switch len(parts) {
	case 0:
		return fmt.Errorf("%w: empty path", backend.ErrUnexpectedError)
	case 1:
		// Check if there is a file with this name
		if _, ok := d.files[parts[0]]; ok {
			return fmt.Errorf("%w: %q", backend.ErrFileAlreadyExists, path)
		}

		// Check if a directory with this name already exists
		if _, ok := d.directories[parts[0]]; ok {
			return fmt.Errorf("%w: %q", backend.ErrDirectoryAlreadyExists, path)
		}

		d.directories[parts[0]] = newDirectory()
		return nil
	default:
		// Check if there is a file with this name
		if _, ok := d.files[parts[0]]; ok {
			return fmt.Errorf("%w: %q", backend.ErrIsFile, path)
		}

		// Create the directory if it does not exist
		if _, ok := d.directories[parts[0]]; !ok {
			d.directories[parts[0]] = newDirectory()
		}

		// Create the directory in the child
		return d.directories[parts[0]].CreateDirectory(strings.Join(parts[1:], "/"))
	}
}

func (d *directory) IsDirectory(path string) error {
	// Split path
	parts := strings.Split(path, "/")

	switch len(parts) {
	case 0:
		return fmt.Errorf("%w: empty path", backend.ErrUnexpectedError)
	case 1:
		// Check if there is a file with this name
		if _, ok := d.files[parts[0]]; ok {
			return fmt.Errorf("%w: %q", backend.ErrIsFile, path)
		}

		// Check if there is a directory with this name
		if _, ok := d.directories[parts[0]]; !ok {
			return fmt.Errorf("%w: %q", backend.ErrNotFound, path)
		}

		return nil
	default:
		// Check if there is a child with the first part of path
		if _, ok := d.directories[parts[0]]; !ok {
			return fmt.Errorf("%w: %q", backend.ErrNotFound, path)
		}

		// Check in the child
		return d.directories[parts[0]].IsDirectory(strings.Join(parts[1:], "/"))
	}
}

func (d *directory) CreateFile(path string, chunkSize int) error {
	// Split path
	parts := strings.Split(path, "/")

	switch {
	case len(parts) == 1 && parts[0] == "":
		return fmt.Errorf("%w: empty path", backend.ErrUnexpectedError)
	case len(parts) == 1:
		// Check if there is a file with this name
		if _, ok := d.files[parts[0]]; ok {
			return fmt.Errorf("%w: %q", backend.ErrFileAlreadyExists, path)
		}

		// Check if there is a directory with this name
		if _, ok := d.directories[parts[0]]; ok {
			return fmt.Errorf("%w: %q", backend.ErrDirectoryAlreadyExists, path)
		}

		// Create the file if it does not exist
		if _, ok := d.files[parts[0]]; !ok {
			f, err := newFile(chunkSize)
			if err != nil {
				return err
			}
			d.files[parts[0]] = f
		}

		return nil
	default:
		// Check if there is a file with this name
		if _, ok := d.files[parts[0]]; ok {
			return fmt.Errorf("%w: %q", backend.ErrIsFile, path)
		}

		// Create the directory if it does not exist
		if _, ok := d.directories[parts[0]]; !ok {
			d.directories[parts[0]] = newDirectory()
		}

		// Create the file in the child
		return d.directories[parts[0]].CreateDirectory(strings.Join(parts[1:], "/"))
	}
}

func (d *directory) IsFile(path string) error {
	// Split path
	parts := strings.Split(path, "/")

	switch {
	case len(parts) == 1 && parts[0] == "":
		return fmt.Errorf("%w: empty path", backend.ErrUnexpectedError)
	case len(parts) == 1:
		// Check if there is a directory with this name
		if _, ok := d.directories[parts[0]]; ok {
			return fmt.Errorf("%w: %q", backend.ErrIsDirectory, path)
		}

		// Check if there is a file with this name
		if _, ok := d.files[parts[0]]; !ok {
			return fmt.Errorf("%w: %q", backend.ErrNotFound, path)
		}

		return nil
	default:
		// Check if there is a child with the first part of path
		if _, ok := d.directories[parts[0]]; !ok {
			return fmt.Errorf("%w: %q", backend.ErrNotFound, path)
		}

		// Check in the child
		return d.directories[parts[0]].IsFile(strings.Join(parts[1:], "/"))
	}
}

func (d *directory) ListFiles(path string) ([]string, error) {
	// Split path
	parts := strings.Split(path, "/")

	switch {
	case len(parts) == 1 && parts[0] == "":
		return slices.Collect(maps.Keys(d.files)), nil
	default:
		// Check if there is a file with this name
		if _, ok := d.files[parts[0]]; ok {
			return nil, fmt.Errorf("%w: %q", backend.ErrIsFile, path)
		}

		// Check if there is a directory with this name
		if _, ok := d.directories[parts[0]]; !ok {
			return nil, fmt.Errorf("%w: %q", backend.ErrNotFound, path)
		}

		// List files in the child
		return d.directories[parts[0]].ListFiles(strings.Join(parts[1:], "/"))
	}
}
