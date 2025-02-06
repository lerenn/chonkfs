package test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/storage/backend"
)

func (suite *BackEndSuite) TestCreateDirectory() {
	err := suite.BackEnd.CreateDirectory(context.Background(), "toto")
	suite.NoError(err)

	err = suite.BackEnd.IsDirectory(context.Background(), "toto")
	suite.NoError(err)
}

func (suite *BackEndSuite) TestCreateDirectoryWithExtraSlash() {
	err := suite.BackEnd.CreateDirectory(context.Background(), "/toto/")
	suite.NoError(err)

	err = suite.BackEnd.IsDirectory(context.Background(), "toto")
	suite.NoError(err)
}

func (suite *BackEndSuite) TestCreateDirectoryOnExistingFile() {
	err := suite.BackEnd.CreateFile(context.Background(), "toto", 4096)
	suite.NoError(err)

	err = suite.BackEnd.CreateDirectory(context.Background(), "toto")
	suite.ErrorIs(err, backend.ErrFileAlreadyExists)
}

func (suite *BackEndSuite) TestCreateDirectoryOnExistingDirectory() {
	err := suite.BackEnd.CreateDirectory(context.Background(), "toto")
	suite.NoError(err)

	err = suite.BackEnd.CreateDirectory(context.Background(), "toto")
	suite.ErrorIs(err, backend.ErrDirectoryAlreadyExists)
}

func (suite *BackEndSuite) TestIsDirectoryWhenIsFile() {
	err := suite.BackEnd.CreateFile(context.Background(), "toto", 4096)
	suite.NoError(err)

	err = suite.BackEnd.IsDirectory(context.Background(), "toto")
	suite.ErrorIs(err, backend.ErrIsFile)
}

func (suite *BackEndSuite) TestListFiles() {
	err := suite.BackEnd.CreateFile(context.Background(), "1", 4096)
	suite.NoError(err)
	err = suite.BackEnd.CreateFile(context.Background(), "2", 4096)
	suite.NoError(err)
	err = suite.BackEnd.CreateFile(context.Background(), "3", 4096)
	suite.NoError(err)

	err = suite.BackEnd.CreateDirectory(context.Background(), "dir")
	suite.NoError(err)

	files, err := suite.BackEnd.ListFiles(context.Background(), "/")
	suite.NoError(err)
	suite.ElementsMatch([]string{"1", "2", "3"}, files)
}
