package layer

import (
	"context"
	"testing"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

type FileSuite struct {
	BackEnd    storage.Directory
	Underlayer storage.Directory
	test.FileSuite
}

func (suite *FileSuite) SetupTest() {
	suite.BackEnd = mem.NewDirectory()
	suite.Underlayer = mem.NewDirectory()
	suite.Directory = NewDirectory(suite.BackEnd, suite.Underlayer)
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

	// Check the file on the backend
	_, err = suite.BackEnd.GetFile(context.Background(), "FileA")
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

	// Check the file on the backend
	dfile, err := suite.BackEnd.GetFile(context.Background(), "FileA")
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
