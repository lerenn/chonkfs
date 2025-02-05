package test

import "github.com/lerenn/chonkfs/pkg/storage/backend"

func (suite *BackEndSuite) TestCreateFile() {
	err := suite.BackEnd.CreateFile(nil, "toto", 4096)
	suite.NoError(err)

	err = suite.BackEnd.IsFile(nil, "toto")
	suite.NoError(err)
}

func (suite *BackEndSuite) TestCreateFileWithExtraSlash() {
	err := suite.BackEnd.CreateFile(nil, "/toto", 4096)
	suite.NoError(err)

	err = suite.BackEnd.IsFile(nil, "toto")
	suite.NoError(err)
}

func (suite *BackEndSuite) TestCreateFileOnExistingFile() {
	err := suite.BackEnd.CreateFile(nil, "toto", 4096)
	suite.NoError(err)

	err = suite.BackEnd.CreateFile(nil, "toto", 4096)
	suite.ErrorIs(err, backend.ErrFileAlreadyExists)
}

func (suite *BackEndSuite) TestCreateFileOnExistingDirectory() {
	err := suite.BackEnd.CreateDirectory(nil, "toto")
	suite.NoError(err)

	err = suite.BackEnd.CreateFile(nil, "toto", 4096)
	suite.ErrorIs(err, backend.ErrDirectoryAlreadyExists)
}

func (suite *BackEndSuite) TestCreateFileWithZeroChunkSize() {
	err := suite.BackEnd.CreateFile(nil, "toto", 0)
	suite.ErrorIs(err, backend.ErrInvalidChunkSize)
}

func (suite *BackEndSuite) TestIsFileWhenIsDirectory() {
	err := suite.BackEnd.CreateDirectory(nil, "toto")
	suite.NoError(err)

	err = suite.BackEnd.IsFile(nil, "toto")
	suite.ErrorIs(err, backend.ErrIsDirectory)
}
