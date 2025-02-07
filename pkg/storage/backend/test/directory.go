package test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/storage/backend"
	"github.com/stretchr/testify/suite"
)

type DirectorySuite struct {
	Directory backend.Directory
	suite.Suite
}

func (suite *DirectorySuite) TestCreateDirectory() {
	d, err := suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().NoError(err)
	suite.Require().NotNil(d)

	rd, err := suite.Directory.GetDirectory(context.Background(), "toto")
	suite.Require().NoError(err)
	suite.Require().NotNil(rd)
}

func (suite *DirectorySuite) TestCreateDirectoryOnExistingFile() {
	_, err := suite.Directory.CreateFile(context.Background(), "toto", 4096)
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, backend.ErrFileAlreadyExists)
}

func (suite *DirectorySuite) TestCreateDirectoryOnExistingDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, backend.ErrDirectoryAlreadyExists)
}

func (suite *DirectorySuite) TestGetDirectoryWhenDoesNotExist() {
	_, err := suite.Directory.GetDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, backend.ErrNotFound)
}

func (suite *DirectorySuite) TestGetDirectoryWhenIsFile() {
	_, err := suite.Directory.CreateFile(context.Background(), "toto", 4096)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, backend.ErrIsFile)
}

func (suite *DirectorySuite) TestListFiles() {
	_, err := suite.Directory.CreateFile(context.Background(), "1", 4096)
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "2", 4096)
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "3", 4096)
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	files, err := suite.Directory.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Contains(files, "1")
	suite.Require().Contains(files, "2")
	suite.Require().Contains(files, "3")
}

func (suite *DirectorySuite) TestRemoveDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	err = suite.Directory.RemoveDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	_, err = suite.Directory.GetDirectory(context.Background(), "dir")
	suite.Require().ErrorIs(err, backend.ErrNotFound)
}

func (suite *DirectorySuite) TestRemoveDirectoryWhenDoesNotExist() {
	err := suite.Directory.RemoveDirectory(context.Background(), "dir")
	suite.Require().ErrorIs(err, backend.ErrNotFound)
}
