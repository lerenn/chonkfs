package test

import "github.com/lerenn/chonkfs/pkg/storage/backend"

func (suite *BackEndSuite) TestCreateDirectory() {
	err := suite.BackEnd.CreateDirectory(nil, "toto")
	suite.NoError(err)

	err = suite.BackEnd.IsDirectory(nil, "toto")
	suite.NoError(err)
}

func (suite *BackEndSuite) TestCreateDirectoryWithExtraSlash() {
	err := suite.BackEnd.CreateDirectory(nil, "/toto/")
	suite.NoError(err)

	err = suite.BackEnd.IsDirectory(nil, "toto")
	suite.NoError(err)
}

func (suite *BackEndSuite) TestCreateDirectoryOnExistingFile() {
	err := suite.BackEnd.CreateFile(nil, "toto", 4096)
	suite.NoError(err)

	err = suite.BackEnd.CreateDirectory(nil, "toto")
	suite.ErrorIs(err, backend.ErrFileAlreadyExists)
}

func (suite *BackEndSuite) TestCreateDirectoryOnExistingDirectory() {
	err := suite.BackEnd.CreateDirectory(nil, "toto")
	suite.NoError(err)

	err = suite.BackEnd.CreateDirectory(nil, "toto")
	suite.ErrorIs(err, backend.ErrDirectoryAlreadyExists)
}

func (suite *BackEndSuite) TestIsDirectoryWhenIsFile() {
	err := suite.BackEnd.CreateFile(nil, "toto", 4096)
	suite.NoError(err)

	err = suite.BackEnd.IsDirectory(nil, "toto")
	suite.ErrorIs(err, backend.ErrIsFile)
}
