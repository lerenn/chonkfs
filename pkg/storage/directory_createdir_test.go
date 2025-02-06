package storage_test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/storage"
)

func (suite *DirectorySuite) TestCreateDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Check it exists on directory backend
	_, err = suite.DirectoryBackEnd.GetDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Check it exists on underlayer backend
	_, err = suite.UnderlayerBackEnd.GetDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)
}

func (suite *DirectorySuite) TestCreateDirectoryWhenDirectoryAlreadyExists() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Create the same directory again
	_, err = suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}

func (suite *DirectorySuite) TestCreateDirectoryWhenFileAlreadyExists() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "test", 4096)
	suite.Require().NoError(err)

	// Create a directory with the same name
	_, err = suite.Directory.CreateDirectory(context.Background(), "test")
	suite.Require().ErrorIs(err, storage.ErrFileAlreadyExists)
}
