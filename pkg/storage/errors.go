package storage

import "fmt"

var (
	// ErrStorage regroups errors from storage.
	ErrStorage = fmt.Errorf("storage error")
	// ErrDirectoryNotExists happens when the requested directory doesn't exist.
	ErrDirectoryNotExists = fmt.Errorf("%w: directory does not exist", ErrStorage)
	// ErrDirectoryAlreadyExists happens when an already existing directory is making the operation fails.
	ErrDirectoryAlreadyExists = fmt.Errorf("%w: directory already exists", ErrStorage)
	// ErrFileNotExists happens when the requested file doesn't exist.
	ErrFileNotExists = fmt.Errorf("%w: file does not exist", ErrStorage)
	// ErrFileAlreadyExists happens when an already existing file is making the operation fails.
	ErrFileAlreadyExists = fmt.Errorf("%w: file already exists", ErrStorage)
	// ErrInvalidChunkNb happens when the chunk number is invalid.
	ErrInvalidChunkNb = fmt.Errorf("%w: invalid chunk number", ErrStorage)
	// ErrInvalidStartOffset happens when the start offset is invalid.
	ErrInvalidStartOffset = fmt.Errorf("%w: invalid start offset", ErrStorage)
	// ErrInvalidEndOffset happens when the end offset is invalid.
	ErrInvalidEndOffset = fmt.Errorf("%w: invalid end offset", ErrStorage)
	// ErrInvalidChunkSize happens when the chunk size is invalid.
	ErrInvalidChunkSize = fmt.Errorf("%w: invalid chunk size", ErrStorage)
	// ErrNoChunk happens when there is no chunk in the file.
	ErrNoChunk = fmt.Errorf("%w: no chunk in file", ErrStorage)
	// ErrLastChunkNotFull happens when the last chunk is not full.
	ErrLastChunkNotFull = fmt.Errorf("%w: last chunk is not full", ErrStorage)
)
