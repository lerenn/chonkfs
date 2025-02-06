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
	err := suite.Directory.CreateFile(nil, "toto", 4096)
	suite.NoError(err)

	err = suite.Directory.IsFile(nil, "toto")
	suite.NoError(err)
}

func (suite *FileSuite) TestCreateFileOnExistingFile() {
	err := suite.Directory.CreateFile(nil, "toto", 4096)
	suite.NoError(err)

	err = suite.Directory.CreateFile(nil, "toto", 4096)
	suite.ErrorIs(err, backend.ErrFileAlreadyExists)
}

func (suite *FileSuite) TestCreateFileOnExistingDirectory() {
	err := suite.Directory.CreateDirectory(nil, "toto")
	suite.NoError(err)

	err = suite.Directory.CreateFile(nil, "toto", 4096)
	suite.ErrorIs(err, backend.ErrDirectoryAlreadyExists)
}

func (suite *FileSuite) TestCreateFileWithZeroChunkSize() {
	err := suite.Directory.CreateFile(nil, "toto", 0)
	suite.ErrorIs(err, backend.ErrInvalidChunkSize)
}

func (suite *FileSuite) TestIsFileWhenIsDirectory() {
	err := suite.Directory.CreateDirectory(nil, "toto")
	suite.NoError(err)

	err = suite.Directory.IsFile(nil, "toto")
	suite.ErrorIs(err, backend.ErrIsDirectory)
}
