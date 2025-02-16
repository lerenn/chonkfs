package disk

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
)

const (
	metadataFileName = ".metadata"
)

type file struct {
	path string
}

func newFile(path string, info info.File) (*file, error) {
	if info.ChunkSize <= 0 {
		return nil, fmt.Errorf("%w: chunk size must be greater than 0", storage.ErrInvalidChunkSize)
	}

	return &file{
		path: path,
	}, nil
}

func writeMetadata(p string, info info.File) error {
	metadataPath := path.Join(p, metadataFileName)

	// Encode info
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	// Create metadata file
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return err
	}

	return nil
}

func readMetadata(p string) (info.File, error) {
	metadataPath := path.Join(p, metadataFileName)

	// Read metadata
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return info.File{}, err
	}

	var info info.File
	if err := json.Unmarshal(data, &info); err != nil {
		return info, err
	}

	return info, nil
}

func getChunkName(nb int) string {
	return fmt.Sprintf("chunk-%d.dat", nb)
}

func (f *file) getChunckPath(nb int) string {
	return path.Join(f.path, getChunkName(nb))
}

func (f *file) checkImportChunkParams(info info.File, index int, data []byte) error {
	// Check if chunk index is correct
	if index < 0 || index >= info.ChunksCount {
		return fmt.Errorf("%w: %d", storage.ErrInvalidChunkNb, index)
	}

	// Check if the chunk is empty
	chunkPath := f.getChunckPath(index)
	if _, err := os.Stat(chunkPath); err == nil {
		return fmt.Errorf("%w: %d", storage.ErrChunkAlreadyExists, index)
	} else if !os.IsNotExist(err) {
		return err
	}

	// Check if length of data is correct
	if (len(data) != info.ChunkSize && index != info.ChunksCount-1) || len(data) > info.ChunkSize {
		return fmt.Errorf("%w: %d", storage.ErrInvalidChunkSize, len(data))
	}

	return nil
}

// ImportChunk imports a chunk of data.
func (f *file) ImportChunk(ctx context.Context, index int, data []byte) error {
	// Get info
	info, err := f.GetInfo(ctx)
	if err != nil {
		return err
	}

	// Check params
	if err := f.checkImportChunkParams(info, index, data); err != nil {
		return err
	}

	// Import data
	chunkPath := f.getChunckPath(index)
	if err := os.WriteFile(chunkPath, data, 0644); err != nil {
		return err
	}

	// If this is the last chunk, set the last chunk size
	if index == info.ChunksCount-1 {
		info.LastChunkSize = len(data)
		if err := f.saveInfo(info); err != nil {
			return err
		}
	}

	return nil
}

// GetInfo returns the file info.
func (f *file) GetInfo(_ context.Context) (info.File, error) {
	return readMetadata(f.path)
}

func (f *file) saveInfo(info info.File) error {
	return writeMetadata(f.path, info)
}

func (f *file) checkReadWriteChunkParams(info info.File, index int, offset int) error {
	// Check if chunk index is correct
	if index < 0 || index >= info.ChunksCount {
		return fmt.Errorf("%w: %d", storage.ErrInvalidChunkNb, index)
	}

	// Check if there is data to read
	chunkPath := f.getChunckPath(index)
	if _, err := os.Stat(chunkPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w", storage.ErrChunkNotFound)
		}
		return err
	}

	// Check if offset is correct
	if offset < 0 || offset >= info.ChunkSize {
		return fmt.Errorf("%w: %d", storage.ErrInvalidOffset, offset)
	}

	// Check if this is the last chunk, that the offset is correct$
	if index == info.ChunksCount-1 && offset >= info.LastChunkSize {
		return fmt.Errorf("%w: %d", storage.ErrInvalidOffset, offset)
	}

	return nil
}

// WriteChunk writes a chunk of data.
func (f *file) WriteChunk(ctx context.Context, index int, data []byte, offset int) (int, error) {
	// Get info
	info, err := f.GetInfo(ctx)
	if err != nil {
		return 0, err
	}

	// Check params
	if err := f.checkReadWriteChunkParams(info, index, offset); err != nil {
		return 0, err
	}

	// Open file
	chunkPath := f.getChunckPath(index)
	file, err := os.OpenFile(chunkPath, os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}

	// Limit data to write if it is too long
	if len(data) > info.ChunkSize-offset {
		data = data[:info.ChunkSize-offset]
	}

	// Write data
	n, err := file.WriteAt(data, int64(offset))
	if err != nil {
		return 0, err
	}

	// Close file
	return n, file.Close()
}

// ReadChunk reads a chunk of data.
func (f *file) ReadChunk(ctx context.Context, index int, data []byte, offset int) (int, error) {
	// Get info
	info, err := f.GetInfo(ctx)
	if err != nil {
		return 0, err
	}

	// Check params
	if err := f.checkReadWriteChunkParams(info, index, offset); err != nil {
		return 0, err
	}

	// Read data
	chunkPath := f.getChunckPath(index)
	chunkData, err := os.ReadFile(chunkPath)
	if err != nil {
		return 0, err
	}

	return copy(data, chunkData[offset:]), nil
}

// ResizeChunksNb resizes the number of chunks.
func (f *file) ResizeChunksNb(ctx context.Context, size int) error {
	// Check size is correct
	if size < 0 {
		return fmt.Errorf("%w: %d", storage.ErrInvalidChunkNb, size)
	}

	// Get actual info
	info, err := f.GetInfo(ctx)
	if err != nil {
		return err
	}

	// Check the last chunk size is full
	if info.ChunksCount > 0 && info.LastChunkSize != info.ChunkSize {
		return fmt.Errorf("%w", storage.ErrLastChunkNotFull)
	}

	// Resize chunks
	if size > info.ChunksCount {
		// Add chunks
		for i := info.ChunksCount; i < size; i++ {
			path := f.getChunckPath(i)
			if err := os.WriteFile(path, make([]byte, info.ChunkSize), 0644); err != nil {
				return err
			}
		}
	} else {
		// Remove chunks
		for i := size; i < info.ChunksCount; i++ {
			path := f.getChunckPath(i)
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}

	// Update info
	info.ChunksCount = size
	info.LastChunkSize = info.ChunkSize
	return f.saveInfo(info)
}

func (f *file) checkResizeLastChunkParams(info info.File, size int) error {
	// Check size is correct
	if size < 0 || size > info.ChunkSize {
		return fmt.Errorf("%w: %d", storage.ErrInvalidChunkSize, size)
	}

	// Check if there is a last chunk
	if info.ChunksCount == 0 {
		return fmt.Errorf("%w", storage.ErrNoChunk)
	}

	// Check if the last chunk is present
	lastChunkPath := f.getChunckPath(info.ChunksCount - 1)
	if _, err := os.Stat(lastChunkPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w", storage.ErrChunkNotFound)
		}
		return err
	}

	return nil
}

// ResizeLastChunk resizes the last chunk.
func (f *file) ResizeLastChunk(ctx context.Context, size int) (changed int, err error) {
	// Get actual info
	info, err := f.GetInfo(ctx)
	if err != nil {
		return 0, err
	}

	// Check params
	if err := f.checkResizeLastChunkParams(info, size); err != nil {
		return 0, err
	}

	// Resize last chunk
	lastChunkPath := f.getChunckPath(info.ChunksCount - 1)
	lastChunkSize := info.LastChunkSize
	if size > lastChunkSize {
		// Append data
		if err := appendFile(lastChunkPath, make([]byte, size-lastChunkSize)); err != nil {
			return 0, err
		}
	} else {
		// Remove data
		if err := os.Truncate(lastChunkPath, int64(size)); err != nil {
			return 0, err
		}
	}

	// Set size
	info.LastChunkSize = size
	if err := f.saveInfo(info); err != nil {
		return 0, err
	}

	return size - lastChunkSize, nil
}

func appendFile(path string, data []byte) error {
	// Open file
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Write data (append)
	if _, err := f.Write(data); err != nil {
		return err
	}

	// Close file
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
