package layer

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/lerenn/chonkfs/pkg/storage/test"
)

type FileSuite struct {
	Upperlayer storage.Directory
	Underlayer storage.Directory
	test.FileSuite
}

func (suite *FileSuite) TestCreateFileAndCheckUnderlayer() {
	// Create a directory
	_, err := suite.Directory.CreateFile(context.Background(), "FileA", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Check it exists on underlayer
	_, err = suite.Underlayer.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
}

func (suite *FileSuite) TestGetInfoWhenFileExistsOnlyOnUnderlayer() {
	// Create a file on underlayer
	_, err := suite.Underlayer.CreateFile(context.Background(), "FileA", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	suite.Require().NotNil(file)

	// Get the file info
	info, err := file.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(4096, info.ChunkSize)
}

func (suite *FileSuite) TestResizeChunksNbOnBackendAndUnderlayer() {
	// Create a file
	file, err := suite.Directory.CreateFile(context.Background(), "FileA", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Resize the file
	err = file.ResizeChunksNb(context.Background(), 3)
	suite.Require().NoError(err)

	// Check the file on the upperlayer
	_, err = suite.Upperlayer.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	info, err := file.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(3, info.ChunksCount)

	// Check the file on the underlayer
	_, err = suite.Underlayer.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	info, err = file.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(3, info.ChunksCount)
}

func (suite *FileSuite) TestResizeChunksNbOnUnderlayerOnly() {
	// Create a file on underlayer
	_, err := suite.Underlayer.CreateFile(context.Background(), "FileA", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)

	// Resize the file
	err = file.ResizeChunksNb(context.Background(), 3)
	suite.Require().NoError(err)

	// Check the file
	_, err = suite.Directory.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	info, err := file.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(3, info.ChunksCount)

	// Check the file on the underlayer
	_, err = suite.Underlayer.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	info, err = file.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(3, info.ChunksCount)
}

func (suite *FileSuite) TestResizeLastChunkOnBackendAndUnderlayer() {
	// Create a file
	file, err := suite.Directory.CreateFile(context.Background(), "FileA", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Resize the file
	err = file.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)

	// Resize the last chunk
	changed, err := file.ResizeLastChunk(context.Background(), 2048)
	suite.Require().NoError(err)
	suite.Require().Equal(2048-4096, changed)

	// Check the file on the upperlayer
	dfile, err := suite.Upperlayer.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	info, err := dfile.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(2048, info.LastChunkSize)

	// Check the file on the underlayer
	ufile, err := suite.Underlayer.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	info, err = ufile.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(2048, info.LastChunkSize)
}

func (suite *FileSuite) TestResizeLastChunkOnUnderlayerOnly() {
	// Create a file on underlayer
	ufile, err := suite.Underlayer.CreateFile(context.Background(), "FileA", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Resize the file
	err = ufile.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)

	// Resize the last chunk
	changed, err := ufile.ResizeLastChunk(context.Background(), 2048)
	suite.Require().NoError(err)
	suite.Require().Equal(2048-4096, changed)

	// Check the file on layer
	file, err := suite.Directory.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	info, err := file.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(2048, info.LastChunkSize)

	// Check the file on the underlayer
	file, err = suite.Underlayer.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	info, err = file.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(2048, info.LastChunkSize)
}

func (suite *FileSuite) TestResizeLastChunkWhenUnderlayerOnly() {
	// Create a file on underlayer
	ufile, err := suite.Underlayer.CreateFile(context.Background(), "FileA", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Resize the file
	err = ufile.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)

	// Resize the last chunk
	changed, err := file.ResizeLastChunk(context.Background(), 2048)
	suite.Require().NoError(err)
	suite.Require().Equal(2048-4096, changed)

	// Check the file
	dFile, err := suite.Directory.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	info, err := dFile.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(2048, info.LastChunkSize)

	// Check the file on the underlayer
	ufile, err = suite.Underlayer.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	info, err = ufile.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(2048, info.LastChunkSize)
}

func (suite *FileSuite) TestReadChunkWhenUnderlayerOnly() {
	// Create a file on underlayer
	ufile, err := suite.Underlayer.CreateFile(context.Background(), "FileA", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Add a chunk
	err = ufile.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)

	// Write a chunk
	written, err := ufile.WriteChunk(context.Background(), 0, []byte("Hello, World!"), 0)
	suite.Require().NoError(err)
	suite.Require().Equal(13, written)

	// Get file from upper layer
	file, err := suite.Directory.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)

	// Read a chunk from upper layer
	data := make([]byte, 4096)
	read, err := file.ReadChunk(context.Background(), 0, data, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(4096, read)
	suite.Require().Equal("Hello, World!", string(data[:13]))
}

func (suite *FileSuite) TestWriteChunkWhenUnderlayerOnly() {
	// Create a file on underlayer
	ufile, err := suite.Underlayer.CreateFile(context.Background(), "FileA", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Add a chunk
	err = ufile.ResizeChunksNb(context.Background(), 1)
	suite.Require().NoError(err)

	// Get file from upper layer
	file, err := suite.Directory.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)

	// Write a chunk from upper layer
	written, err := file.WriteChunk(context.Background(), 0, []byte("Hello, World!"), 0)
	suite.Require().NoError(err)
	suite.Require().Equal(13, written)

	// Read a chunk from underlayer
	data := make([]byte, 4096)
	read, err := ufile.ReadChunk(context.Background(), 0, data, 0)
	suite.Require().NoError(err)
	suite.Require().Equal(4096, read)
	suite.Require().Equal("Hello, World!", string(data[:13]))
}
