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
	err := suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.NoError(err)

	err = suite.Directory.IsDirectory(context.Background(), "toto")
	suite.NoError(err)
}

func (suite *DirectorySuite) TestCreateDirectoryOnExistingFile() {
	err := suite.Directory.CreateFile(context.Background(), "toto", 4096)
	suite.NoError(err)

	err = suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.ErrorIs(err, backend.ErrFileAlreadyExists)
}

func (suite *DirectorySuite) TestCreateDirectoryOnExistingDirectory() {
	err := suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.NoError(err)

	err = suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.ErrorIs(err, backend.ErrDirectoryAlreadyExists)
}

func (suite *DirectorySuite) TestIsDirectoryWhenIsFile() {
	err := suite.Directory.CreateFile(context.Background(), "toto", 4096)
	suite.NoError(err)

	err = suite.Directory.IsDirectory(context.Background(), "toto")
	suite.ErrorIs(err, backend.ErrIsFile)
}

func (suite *DirectorySuite) TestListFiles() {
	err := suite.Directory.CreateFile(context.Background(), "1", 4096)
	suite.NoError(err)
	err = suite.Directory.CreateFile(context.Background(), "2", 4096)
	suite.NoError(err)
	err = suite.Directory.CreateFile(context.Background(), "3", 4096)
	suite.NoError(err)

	err = suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.NoError(err)

	files, err := suite.Directory.ListFiles(context.Background())
	suite.NoError(err)
	suite.ElementsMatch([]string{"1", "2", "3"}, files)
}
