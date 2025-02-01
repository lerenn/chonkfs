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

// TestResizeChunksAndChunksCount tests resizing the chunks and getting the chunks count.
func (suite *FileSuite) TestResizeChunksAndChunksCount() {
	// Create a file from the directory
	file, err := suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().NoError(err)

	// Get the underlaying file
	ufile, err := suite.Underlayer.GetFile(context.Background(), "file")
	suite.Require().NoError(err)

	// Check the chunks count
	count, err := file.ChunksCount(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(0, count)

	// Check the chunks count in the underlayer
	ucount, err := ufile.ChunksCount(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(0, ucount)

	// Resize the chunks count
	err = file.ResizeChunksNb(context.Background(), 12)
	suite.Require().NoError(err)

	// Check the chunks count
	count, err = file.ChunksCount(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(12, count)

	// Check the chunks count in the underlayer
	count, err = ufile.ChunksCount(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(12, count)

	// Write last chunk with success
	_, err = file.WriteChunk(context.Background(), 11, 0, nil, []byte("x"))
	suite.Require().NoError(err)

	// Write next chunk that doesn't exists
	_, err = file.WriteChunk(context.Background(), 12, 0, nil, []byte("x"))
	suite.Require().ErrorIs(err, storage.ErrInvalidChunkNb)
}

// TestReadWriteChunk tests reading and writing a chunk.
func (suite *FileSuite) TestReadWriteChunk() {
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
	_, err = file.WriteChunk(context.Background(), 0, 0, nil, data)
	suite.Require().NoError(err)

	// Read the data
	readData := make([]byte, 13)
	_, err = file.ReadChunk(context.Background(), 0, readData, 0, nil)
	suite.Require().NoError(err)
	suite.Require().Equal(data, readData)

	// Read the data on the underlayer
	uReadData := make([]byte, 13)
	_, err = ufile.ReadChunk(context.Background(), 0, uReadData, 0, nil)
	suite.Require().NoError(err)
	suite.Require().Equal(data, uReadData)
}

// TestResizeLastChunkAndLastChunkSize tests resizing the last chunk and getting the last chunk size.
func (suite *FileSuite) TestResizeLastChunkAndLastChunkSize() {
	// Create a file from the directory
	file, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	// Get the underlaying file
	ufile, err := suite.Underlayer.GetFile(context.Background(), "file")
	suite.Require().NoError(err)

	// Resize the last chunk, but fails because no chunk
	_, err = file.ResizeLastChunk(context.Background(), 12)
	suite.Require().ErrorIs(err, storage.ErrNoChunk)

	// Add a chunk
	err = file.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)

	// Check the last chunk size
	size, err := file.LastChunkSize(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(4096, size)

	// Check the last chunk size on the underlayer
	usize, err := ufile.LastChunkSize(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(4096, usize)

	// Resize the last chunk and returns the subtract values
	changed, err := file.ResizeLastChunk(context.Background(), 12)
	suite.Require().NoError(err)
	suite.Require().Equal(12-4096, changed)

	// Check the last chunk size
	size, err = file.LastChunkSize(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(12, size)

	// Check the last chunk size in the underlayer
	usize, err = ufile.LastChunkSize(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(12, usize)

	// Resize the last chunk and returns the subtract values
	changed, err = file.ResizeLastChunk(context.Background(), 24)
	suite.Require().NoError(err)
	suite.Require().Equal(12, changed)

	// Check the last chunk size
	size, err = file.LastChunkSize(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(24, size)

	// Check the last chunk size in the underlayer
	usize, err = ufile.LastChunkSize(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(24, usize)
}

// TestInfo tests getting the file info.
func (suite *FileSuite) TestInfo() {
	// Create a file from the directory
	file, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	// Check the info
	info, err := file.Info(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(storage.FileInfo{ChunkSize: 4096}, info)
}

// TestSize tests getting the file size.
func (suite *FileSuite) TestSize() {
	// Create a file from the directory
	file, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	// Check the size
	size, err := file.Size(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(0, size)

	// Create a new chunk and resize it to contain the data
	err = file.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)
	changed, err := file.ResizeLastChunk(context.Background(), 13)
	suite.Require().NoError(err)
	suite.Require().Equal(13-4096, changed)

	// Write data
	data := []byte("Hello, World!")
	_, err = file.WriteChunk(context.Background(), 0, 0, nil, data)
	suite.Require().NoError(err)

	// Check the size
	size, err = file.Size(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(13, size)
}
