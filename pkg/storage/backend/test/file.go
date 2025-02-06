package test

import (
	"github.com/lerenn/chonkfs/pkg/storage/backend"
	"github.com/stretchr/testify/suite"
)

type FileSuite struct {
	Directory backend.Directory
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
	suite.Require().ErrorIs(err, backend.ErrFileAlreadyExists)
}

func (suite *FileSuite) TestCreateFileOnExistingDirectory() {
	_, err := suite.Directory.CreateDirectory(nil, "toto")
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateFile(nil, "toto", 4096)
	suite.Require().ErrorIs(err, backend.ErrDirectoryAlreadyExists)
}

func (suite *FileSuite) TestCreateFileWithZeroChunkSize() {
	_, err := suite.Directory.CreateFile(nil, "toto", 0)
	suite.Require().ErrorIs(err, backend.ErrInvalidChunkSize)
}

func (suite *FileSuite) TestIsFileWhenIsDirectory() {
	_, err := suite.Directory.CreateDirectory(nil, "toto")
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(nil, "toto")
	suite.Require().ErrorIs(err, backend.ErrIsDirectory)
}

func (suite *FileSuite) TestGetInfoFromEmptyFile() {
	f, err := suite.Directory.CreateFile(nil, "toto", 4096)
	suite.Require().NoError(err)

	info, err := f.GetInfo(nil)
	suite.Require().NoError(err)

	suite.EqualValues(0, info.ChunksCount)
}
