package test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/stretchr/testify/suite"
)

type FileSuite struct {
	Directory storage.Directory
	suite.Suite
}

func (suite *FileSuite) TestCreateFile() {
	f, err := suite.Directory.CreateFile(context.Background(), "toto", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(f)

	rf, err := suite.Directory.GetFile(context.Background(), "toto")
	suite.Require().NoError(err)
	suite.Require().NotNil(rf)
}

func (suite *FileSuite) TestCreateFileOnExistingFile() {
	_, err := suite.Directory.CreateFile(context.Background(), "toto", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateFile(context.Background(), "toto", info.File{
		ChunkSize: 4096,
	})
	suite.Require().ErrorIs(err, storage.ErrFileAlreadyExists)
}

func (suite *FileSuite) TestCreateFileOnExistingDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateFile(context.Background(), "toto", info.File{
		ChunkSize: 4096,
	})
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}

func (suite *FileSuite) TestCreateFileWithZeroChunkSize() {
	_, err := suite.Directory.CreateFile(context.Background(), "toto", info.File{
		ChunkSize: 0,
	})
	suite.Require().ErrorIs(err, storage.ErrInvalidChunkSize)
}

func (suite *FileSuite) TestIsFileWhenIsDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(context.Background(), "toto")
	suite.Require().ErrorIs(err, storage.ErrIsDirectory)
}

func (suite *FileSuite) TestGetInfoFromEmptyFile() {
	f, err := suite.Directory.CreateFile(context.Background(), "toto", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	info, err := f.GetInfo(context.Background())
	suite.Require().NoError(err)

	suite.Require().Equal(0, info.ChunksCount)
	suite.Require().Equal(4096, info.ChunkSize)
}

func (suite *FileSuite) TestResizeChunksNb() {
	f, err := suite.Directory.CreateFile(context.Background(), "toto", info.File{
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

func (suite *FileSuite) TestResizeLastChunk() {
	f, err := suite.Directory.CreateFile(context.Background(), "toto", info.File{
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

func (suite *FileSuite) TestResizeLastChunkWithInvalidSizes() {
	f, err := suite.Directory.CreateFile(context.Background(), "toto", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	_, err = f.ResizeLastChunk(context.Background(), -1)
	suite.Require().ErrorIs(err, storage.ErrInvalidChunkSize)

	_, err = f.ResizeLastChunk(context.Background(), 4097)
	suite.Require().ErrorIs(err, storage.ErrInvalidChunkSize)
}
