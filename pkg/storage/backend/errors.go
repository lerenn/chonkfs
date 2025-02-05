package backend

import "fmt"

var (
	// ErrBackend is the error returned by the backend
	ErrBackend = fmt.Errorf("backend error")
	// ErrNotFound is the error returned when the backend does not find the requested resource
	ErrNotFound = fmt.Errorf("%w: not found", ErrBackend)
	// ErrIsFile is the error returned when the backend finds a file instead of a directory
	ErrIsFile = fmt.Errorf("%w: is a file", ErrBackend)
	// ErrIsDirectory is the error returned when the backend finds a directory instead of a file
	ErrIsDirectory = fmt.Errorf("%w: is a directory", ErrBackend)
	// ErrUnexpectedError is the error returned when the backend encounters an unexpected error
	ErrUnexpectedError = fmt.Errorf("%w: unexpected error", ErrBackend)
	// ErrFileAlreadyExists is the error returned when the backend finds a file with the same name
	ErrFileAlreadyExists = fmt.Errorf("%w: file already exists", ErrBackend)
	// ErrDirectoryAlreadyExists is the error returned when the backend finds a directory with the same name
	ErrDirectoryAlreadyExists = fmt.Errorf("%w: directory already exists", ErrBackend)
	// ErrInvalidChunkSize is the error returned when the chunk size is invalid
	ErrInvalidChunkSize = fmt.Errorf("%w: invalid chunk size", ErrBackend)
)
