package storage

import "fmt"

var (
	ErrStorage                = fmt.Errorf("storage error")
	ErrDirectoryNotExists     = fmt.Errorf("%w: directory does not exist", ErrStorage)
	ErrDirectoryAlreadyExists = fmt.Errorf("%w: directory already exists", ErrStorage)
	ErrFileNotExists          = fmt.Errorf("%w: file does not exist", ErrStorage)
	ErrFileAlreadyExists      = fmt.Errorf("%w: file already exists", ErrStorage)
	ErrInvalidChunkNb         = fmt.Errorf("%w: invalid chunk number", ErrStorage)
	ErrInvalidStartOffset     = fmt.Errorf("%w: invalid start offset", ErrStorage)
	ErrInvalidEndOffset       = fmt.Errorf("%w: invalid end offset", ErrStorage)
	ErrInvalidChunkSize       = fmt.Errorf("%w: invalid chunk size", ErrStorage)
	ErrNoChunk                = fmt.Errorf("%w: no chunk in file", ErrStorage)
	ErrLastChunkNotFull       = fmt.Errorf("%w: last chunk is not full", ErrStorage)
)
