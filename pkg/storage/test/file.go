package test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/stretchr/testify/suite"
)

// FileSuite is a test suite for a file.
type FileSuite struct {
	Underlayer storage.Directory
	Directory  storage.Directory
	suite.Suite
}

// TestResizeChunks tests resizing the chunks and getting the chunks count.
func (suite *FileSuite) TestResizeChunks() {
	// Create a file from the directory
	file, err := suite.Directory.CreateFile(context.Background(), "file", 2)
	suite.Require().NoError(err)

	// Check the chunks count
	info, err := file.Info(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(0, info.ChunksCount)

	// Resize the chunks count
	err = file.ResizeChunksNb(context.Background(), 12)
	suite.Require().NoError(err)

	// Check the chunks count
	info, err = file.Info(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(12, info.ChunksCount)

	// Write last chunk with success
	_, err = file.WriteChunk(context.Background(), 11, []byte("x"), 0)
	suite.Require().NoError(err)

	// Write next chunk that doesn't exists
	_, err = file.WriteChunk(context.Background(), 12, []byte("x"), 0)
	suite.Require().ErrorIs(err, storage.ErrInvalidChunkNb)
}

// TestResizeChunksInUnderlayer tests resizing the chunks and
// getting the chunks count is passed to underlayer.
func (suite *FileSuite) TestResizeChunksInUnderlayer() {
	// Create a file from the directory
	file, err := suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().NoError(err)

	// Get the underlaying file
	ufile, err := suite.Underlayer.GetFile(context.Background(), "file")
	suite.Require().NoError(err)

	// Check the chunks count in the underlayer
	uinfo, err := ufile.Info(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(0, uinfo.ChunksCount)

	// Resize the chunks count
	err = file.ResizeChunksNb(context.Background(), 12)
	suite.Require().NoError(err)

	// Check the chunks count in the underlayer
	uinfo, err = ufile.Info(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(12, uinfo.ChunksCount)
}

// TestResizeChunksWhenOnlyUnderlayerExists tests resizing the chunks and
// getting the chunks count is passed to underlayer when the file only exists in the underlayer.
func (suite *FileSuite) TestResizeChunksWhenOnlyUnderlayerExists() {
	// Create a file in the underlayer
	ufile, err := suite.Underlayer.CreateFile(context.Background(), "file", 1)
	suite.Require().NoError(err)

	// Resize the chunks count in the underlayer
	err = ufile.ResizeChunksNb(context.Background(), 12)
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "file")
	suite.Require().NoError(err)

	// Check the chunks count in the underlayer
	info, err := file.Info(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(12, info.ChunksCount)
}

// TestReadWriteChunk tests reading and writing a chunk.
func (suite *FileSuite) TestReadWriteChunk() {
	// Create a file from the directory
	file, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	// Create a new chunk and resize it to contain the data
	err = file.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)
	changed, err := file.ResizeLastChunk(context.Background(), 13)
	suite.Require().NoError(err)
	suite.Require().Equal(13-4096, changed)

	// Write data
	data := []byte("Hello, World!")
	_, err = file.WriteChunk(context.Background(), 0, data, 0)
	suite.Require().NoError(err)

	// Read the data
	readData := make([]byte, 13)
	_, err = file.ReadChunk(context.Background(), 0, readData, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(data, readData)
}

// TestReadWriteChunkInUnderlayer tests reading and writing a chunk is passed to underlayer.
func (suite *FileSuite) TestReadWriteChunkInUnderlayer() {
	// Create a file from the directory
	file, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	// Get the underlaying file
	ufile, err := suite.Underlayer.GetFile(context.Background(), "file")
	suite.Require().NoError(err)

	// Create a new chunk and resize it to contain the data
	err = file.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)
	changed, err := file.ResizeLastChunk(context.Background(), 13)
	suite.Require().NoError(err)
	suite.Require().Equal(13-4096, changed)

	// Write data
	data := []byte("Hello, World!")
	_, err = file.WriteChunk(context.Background(), 0, data, 0)
	suite.Require().NoError(err)

	// Read the data on the underlayer
	uReadData := make([]byte, 13)
	_, err = ufile.ReadChunk(context.Background(), 0, uReadData, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(data, uReadData)
}

// TestReadWriteChunkWhenOnlyUnderlayerExists tests reading and writing a chunk
// is passed to underlayer when the file only exists in the underlayer.
// func (suite *FileSuite) TestReadWriteChunkWhenOnlyUnderlayerExists() {
// 	// Create a file in the underlayer
// 	ufile, err := suite.Underlayer.CreateFile(context.Background(), "file", 4096)
// 	suite.Require().NoError(err)

// 	// Create a new chunk and resize it to contain the data
// 	err = ufile.ResizeChunksNb(context.Background(), 1)
// 	suite.Require().NoError(err)
// 	changed, err := ufile.ResizeLastChunk(context.Background(), 13)
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(13-4096, changed)

// 	// Write data in the underlayer
// 	data := []byte("Hello, World!")
// 	_, err = ufile.WriteChunk(context.Background(), 0, data, 0)
// 	suite.Require().NoError(err)

// 	// Get the file
// 	file, err := suite.Directory.GetFile(context.Background(), "file")
// 	suite.Require().NoError(err)

// 	// Read the data
// 	readData := make([]byte, 13)
// 	_, err = file.ReadChunk(context.Background(), 0, readData, 0)
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(data, readData)
// }

// // TestResizeLastChunkAndLastChunkSize tests resizing the last chunk and getting the last chunk size.
// func (suite *FileSuite) TestResizeLastChunkAndLastChunkSize() {
// 	// Create a file from the directory
// 	file, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
// 	suite.Require().NoError(err)

// 	// Resize the last chunk, but fails because no chunk
// 	_, err = file.ResizeLastChunk(context.Background(), 12)
// 	suite.Require().ErrorIs(err, storage.ErrNoChunk)

// 	// Add a chunk
// 	err = file.ResizeChunksNb(context.Background(), 1)
// 	suite.Require().NoError(err)

// 	// Check the last chunk size
// 	size, err := file.LastChunkSize(context.Background())
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(4096, size)

// 	// Resize the last chunk and returns the subtract values
// 	changed, err := file.ResizeLastChunk(context.Background(), 12)
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(12-4096, changed)

// 	// Check the last chunk size
// 	size, err = file.LastChunkSize(context.Background())
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(12, size)

// 	// Resize the last chunk and returns the subtract values
// 	changed, err = file.ResizeLastChunk(context.Background(), 24)
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(12, changed)

// 	// Check the last chunk size
// 	size, err = file.LastChunkSize(context.Background())
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(24, size)
// }

// // TestResizeLastChunkAndLastChunkSizeInUnderlayer tests resizing the last chunk
// // and getting the last chunk size is passed to underlayer.
// func (suite *FileSuite) TestResizeLastChunkAndLastChunkSizeInUnderlayer() {
// 	// Create a file from the directory
// 	file, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
// 	suite.Require().NoError(err)

// 	// Get the underlaying file
// 	ufile, err := suite.Underlayer.GetFile(context.Background(), "file")
// 	suite.Require().NoError(err)

// 	// Resize the last chunk, but fails because no chunk
// 	_, err = file.ResizeLastChunk(context.Background(), 12)
// 	suite.Require().ErrorIs(err, storage.ErrNoChunk)

// 	// Add a chunk
// 	err = file.ResizeChunksNb(context.Background(), 1)
// 	suite.Require().NoError(err)

// 	// Check the last chunk size on the underlayer
// 	usize, err := ufile.LastChunkSize(context.Background())
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(4096, usize)

// 	// Resize the last chunk and returns the subtract values
// 	changed, err := file.ResizeLastChunk(context.Background(), 12)
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(12-4096, changed)

// 	// Check the last chunk size in the underlayer
// 	usize, err = ufile.LastChunkSize(context.Background())
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(12, usize)

// 	// Resize the last chunk and returns the subtract values
// 	changed, err = file.ResizeLastChunk(context.Background(), 24)
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(12, changed)

// 	// Check the last chunk size in the underlayer
// 	usize, err = ufile.LastChunkSize(context.Background())
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(24, usize)
// }

// TestResizeChunkWhenLastChunkNotFull tests resizing the last chunk when it's not full.
// func (suite *FileSuite) TestResizeChunkWhenLastChunkNotFull() {
// 	// Create a file from the directory
// 	file, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
// 	suite.Require().NoError(err)

// 	// Add a chunk
// 	err = file.ResizeChunksNb(context.Background(), 1)
// 	suite.Require().NoError(err)

// 	// Resize the last chunk and fails because it's not full
// 	_, err = file.ResizeLastChunk(context.Background(), 12)
// 	suite.Require().ErrorIs(err, storage.ErrLastChunkNotFull)
// }

// // TestInfo tests getting the file info.
// func (suite *FileSuite) TestInfo() {
// 	// Create a file from the directory
// 	file, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
// 	suite.Require().NoError(err)

// 	// Check the info
// 	info, err := file.Info(context.Background())
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(storage.FileInfo{ChunkSize: 4096}, info)
// }

// TestInfoInUnderlayer tests getting the file info is passed to underlayer.
// func (suite *FileSuite) TestInfoInUnderlayer() {
// 	// Create a file from the directory
// 	_, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
// 	suite.Require().NoError(err)

// 	// Get the underlaying file
// 	ufile, err := suite.Underlayer.GetFile(context.Background(), "file")
// 	suite.Require().NoError(err)

// 	// Check the info in the underlayer
// 	uinfo, err := ufile.Info(context.Background())
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(storage.FileInfo{ChunkSize: 4096}, uinfo)
// }

// // TestInfoWhenOnlyUnderlayerExists tests getting the file info is passed to underlayer
// // when the file only exists in the underlayer.
// func (suite *FileSuite) TestInfoWhenOnlyUnderlayerExists() {
// 	// Create a file in the underlayer
// 	ufile, err := suite.Underlayer.CreateFile(context.Background(), "file", 4096)
// 	suite.Require().NoError(err)
//
// 	// Get the file
// 	file, err := suite.Directory.GetFile(context.Background(), "file")
// 	suite.Require().NoError(err)
//
// 	// Check the info
// 	info, err := file.Info(context.Background())
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(storage.FileInfo{ChunkSize: 4096}, info)
// }

// // TestSize tests getting the file size.
// func (suite *FileSuite) TestSize() {
// 	// Create a file from the directory
// 	file, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
// 	suite.Require().NoError(err)

// 	// Check the size
// 	size, err := file.Size(context.Background())
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(0, size)

// 	// Create a new chunk and resize it to contain the data
// 	err = file.ResizeChunksNb(context.Background(), 1)
// 	suite.Require().NoError(err)
// 	changed, err := file.ResizeLastChunk(context.Background(), 13)
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(13-4096, changed)

// 	// Write data
// 	data := []byte("Hello, World!")
// 	_, err = file.WriteChunk(context.Background(), 0, 0, nil, data)
// 	suite.Require().NoError(err)

// 	// Check the size
// 	size, err = file.Size(context.Background())
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(13, size)
// }
