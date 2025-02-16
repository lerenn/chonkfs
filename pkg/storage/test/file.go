package test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/stretchr/testify/suite"
)

// FileSuite is a test suite for storage.File implementations.
type FileSuite struct {
	Directory storage.Directory
	suite.Suite
}

// TestCreateFile tests the creation of a file.
func (suite *FileSuite) TestCreateFile() {
	f, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(f)

	rf, err := suite.Directory.GetFile(context.Background(), "file")
	suite.Require().NoError(err)
	suite.Require().NotNil(rf)
}

// TestCreateFileOnExistingFile tests the creation of a file on an existing file.
func (suite *FileSuite) TestCreateFileOnExistingFile() {
	_, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().ErrorIs(err, storage.ErrFileAlreadyExists)
}

// TestCreateFileOnExistingDirectory tests the creation of a file on an existing directory.
func (suite *FileSuite) TestCreateFileOnExistingDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "file")
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}

// TestCreateFileWithZeroChunkSize tests the creation of a file with a zero chunk size.
func (suite *FileSuite) TestCreateFileWithZeroChunkSize() {
	_, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 0,
	})
	suite.Require().ErrorIs(err, storage.ErrInvalidChunkSize)
}

// TestGetFileWhenIsDirectory tests the GetFile method on a directory.
func (suite *FileSuite) TestGetFileWhenIsDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "file")
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().ErrorIs(err, storage.ErrIsDirectory)
}

// TestGetInfoFromEmptyFile tests the GetInfo method on an empty file.
func (suite *FileSuite) TestGetInfoFromEmptyFile() {
	f, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	fInfo, err := f.GetInfo(context.Background())
	suite.Require().NoError(err)

	suite.Require().Equal(info.File{
		Size:          0,
		ChunkSize:     4096,
		ChunksCount:   0,
		LastChunkSize: 0,
	}, fInfo)
}

// TestResizeChunksNb tests the ResizeChunksNb method.
func (suite *FileSuite) TestResizeChunksNb() {
	f, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Resize to superior size
	err = f.ResizeChunksNb(context.Background(), 10)
	suite.Require().NoError(err)

	info, err := f.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(10, info.ChunksCount)

	// Resize to inferior size
	err = f.ResizeChunksNb(context.Background(), 5)
	suite.Require().NoError(err)

	info, err = f.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(5, info.ChunksCount)
}

// TestResizeLastChunk tests the ResizeLastChunk method.
func (suite *FileSuite) TestResizeLastChunk() {
	f, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Resize to have at least one chunk
	err = f.ResizeChunksNb(context.Background(), 4096)
	suite.Require().NoError(err)

	// Resize to 0
	changed, err := f.ResizeLastChunk(context.Background(), 0)
	suite.Require().NoError(err)
	suite.Require().Equal(-4096, changed)

	info, err := f.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(0, info.LastChunkSize)

	// Resize to 4096
	changed, err = f.ResizeLastChunk(context.Background(), 1234)
	suite.Require().NoError(err)
	suite.Require().Equal(1234, changed)

	info, err = f.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(1234, info.LastChunkSize)
}

// TestResizeLastChunkWithInvalidSizes tests the ResizeLastChunk method with invalid
// sizes.
func (suite *FileSuite) TestResizeLastChunkWithInvalidSizes() {
	f, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	_, err = f.ResizeLastChunk(context.Background(), -1)
	suite.Require().ErrorIs(err, storage.ErrInvalidChunkSize)

	_, err = f.ResizeLastChunk(context.Background(), 4097)
	suite.Require().ErrorIs(err, storage.ErrInvalidChunkSize)
}

// TestReadWriteChunk tests the ReadChunk and WriteChunk methods.
func (suite *FileSuite) TestReadWriteChunk() {
	// Create file
	f, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Resize to have at least one chunk
	err = f.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)

	// Write chunk
	buf := []byte("Hello, World!")
	n, err := f.WriteChunk(context.Background(), 0, buf, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(len(buf), n)

	// Read chunk
	rbuf := make([]byte, 4096)
	_, err = f.ReadChunk(context.Background(), 0, rbuf, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(buf, rbuf[:len(buf)])
}

// TestImportChunk tests the ImportChunk method.
func (suite *FileSuite) TestImportChunk() {
	// Create file
	f, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize:   4096,
		ChunksCount: 1,
	})
	suite.Require().NoError(err)

	// Import chunk
	chunk := make([]byte, 4096)
	copy(chunk, []byte("Hello, World!"))
	err = f.ImportChunk(context.Background(), 0, chunk)
	suite.Require().NoError(err)

	// Read chunk
	rbuf := make([]byte, 4096)
	_, err = f.ReadChunk(context.Background(), 0, rbuf, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(chunk[:13], rbuf[:13])
}

// TestImportAlreadyExistingChunk tests the ImportChunk method with an already existing
// chunk.
func (suite *FileSuite) TestImportAlreadyExistingChunk() {
	// Create file
	f, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Resize to have at least one chunk
	err = f.ResizeChunksNb(context.Background(), 4096)
	suite.Require().NoError(err)

	// Import chunk
	chunk := make([]byte, 4096)
	err = f.ImportChunk(context.Background(), 0, chunk)
	suite.Require().ErrorIs(err, storage.ErrChunkAlreadyExists)
}

// TestImportTooBigChunk tests the ImportChunk method with a too big chunk.
func (suite *FileSuite) TestImportTooBigChunk() {
	// Create file
	f, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize:   4096,
		ChunksCount: 1,
	})
	suite.Require().NoError(err)

	// Import chunk
	chunk := make([]byte, 4097)
	err = f.ImportChunk(context.Background(), 0, chunk)
	suite.Require().ErrorIs(err, storage.ErrInvalidChunkSize)
}

// TestReadChunkWithBiggerBufferThanChunk tests the ReadChunk method with a buffer
// bigger than the chunk.
func (suite *FileSuite) TestReadChunkWithBiggerBufferThanChunk() {
	// Create file
	f, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Resize to have at least one chunk
	err = f.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)

	// Write chunk
	buf := []byte("Hello, World!")
	n, err := f.WriteChunk(context.Background(), 0, buf, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(len(buf), n)

	// Read chunk
	rbuf := make([]byte, 8192)
	_, err = f.ReadChunk(context.Background(), 0, rbuf, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(buf, rbuf[:len(buf)])
}
