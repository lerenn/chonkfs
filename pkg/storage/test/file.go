package test

import (
	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/stretchr/testify/suite"
)

type FileSuite struct {
	Directory storage.Directory
	suite.Suite
}

func (suite *FileSuite) TestCreateFile() {
	f, err := suite.Directory.CreateFile(nil, "toto", 4096)
	suite.Require().NoError(err)
	suite.Require().NotNil(f)

	rf, err := suite.Directory.GetFile(nil, "toto")
	suite.Require().NoError(err)
	suite.Require().NotNil(rf)
}

func (suite *FileSuite) TestCreateFileOnExistingFile() {
	_, err := suite.Directory.CreateFile(nil, "toto", 4096)
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateFile(nil, "toto", 4096)
	suite.Require().ErrorIs(err, storage.ErrFileAlreadyExists)
}

func (suite *FileSuite) TestCreateFileOnExistingDirectory() {
	_, err := suite.Directory.CreateDirectory(nil, "toto")
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateFile(nil, "toto", 4096)
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}

func (suite *FileSuite) TestCreateFileWithZeroChunkSize() {
	_, err := suite.Directory.CreateFile(nil, "toto", 0)
	suite.Require().ErrorIs(err, storage.ErrInvalidChunkSize)
}

func (suite *FileSuite) TestIsFileWhenIsDirectory() {
	_, err := suite.Directory.CreateDirectory(nil, "toto")
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(nil, "toto")
	suite.Require().ErrorIs(err, storage.ErrIsDirectory)
}

func (suite *FileSuite) TestGetInfoFromEmptyFile() {
	f, err := suite.Directory.CreateFile(nil, "toto", 4096)
	suite.Require().NoError(err)

	info, err := f.GetInfo(nil)
	suite.Require().NoError(err)

	suite.Require().Equal(0, info.ChunksCount)
	suite.Require().Equal(4096, info.ChunkSize)
}

func (suite *FileSuite) TestResizeChunksNb() {
	f, err := suite.Directory.CreateFile(nil, "toto", 4096)
	suite.Require().NoError(err)

	// Resize to superior size
	err = f.ResizeChunksNb(nil, 10)
	suite.Require().NoError(err)

	info, err := f.GetInfo(nil)
	suite.Require().NoError(err)
	suite.Require().Equal(10, info.ChunksCount)

	// Resize to inferior size
	err = f.ResizeChunksNb(nil, 5)
	suite.Require().NoError(err)

	info, err = f.GetInfo(nil)
	suite.Require().NoError(err)
	suite.Require().Equal(5, info.ChunksCount)
}
