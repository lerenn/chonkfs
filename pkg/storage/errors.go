package storage

import "fmt"

var (
	// ErrStorage regroups errors from storage.
	ErrStorage = fmt.Errorf("storage error")
	// ErrDirectoryNotFound happens when the requested directory doesn't exist.
	ErrDirectoryNotFound = fmt.Errorf("%w: directory not found", ErrStorage)
	// ErrDirectoryAlreadyExists happens when an already existing directory is making the operation fails.
	ErrDirectoryAlreadyExists = fmt.Errorf("%w: directory already exists", ErrStorage)
	// ErrIsDirectory happens when the requested element is a directory.
	ErrIsDirectory = fmt.Errorf("%w: is a directory", ErrStorage)
	// ErrFileNotFound happens when the requested file doesn't exist.
	ErrFileNotFound = fmt.Errorf("%w: file not found", ErrStorage)
	// ErrFileAlreadyExists happens when an already existing file is making the operation fails.
	ErrFileAlreadyExists = fmt.Errorf("%w: file already exists", ErrStorage)
	// ErrIsFile happens when the requested element is a file.
	ErrIsFile = fmt.Errorf("%w: is a file", ErrStorage)
	// ErrInvalidChunkNb happens when the chunk number is invalid.
	ErrInvalidChunkNb = fmt.Errorf("%w: invalid chunk number", ErrStorage)
	// ErrInvalidOffset happens when the offset is invalid.
	ErrInvalidOffset = fmt.Errorf("%w: invalid offset", ErrStorage)
	// ErrRequestTooBig happens when the request is too big.
	ErrRequestTooBig = fmt.Errorf("%w: request too big", ErrStorage)
	// ErrInvalidChunkSize happens when the chunk size is invalid.
	ErrInvalidChunkSize = fmt.Errorf("%w: invalid chunk size", ErrStorage)
	// ErrNoChunk happens when there is no chunk in the file.
	ErrNoChunk = fmt.Errorf("%w: no chunk in file", ErrStorage)
	// ErrLastChunkNotFull happens when the last chunk is not full.
	ErrLastChunkNotFull = fmt.Errorf("%w: last chunk is not full", ErrStorage)
	// ErrChunkNotFound happens when the chunk is not present on the medium, but
	// still valid regarding metadata.
	ErrChunkNotFound = fmt.Errorf("%w: chunk not found", ErrStorage)
	// ErrChunkAlreadyExists happens when the chunk already exists and cannot be imported.
	ErrChunkAlreadyExists = fmt.Errorf("%w: chunk already exists", ErrStorage)
)
