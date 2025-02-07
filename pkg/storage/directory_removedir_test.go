package storage_test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/storage/backend"
)

func (suite *DirectorySuite) TestRemoveDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Remove the directory
	err = suite.Directory.RemoveDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Check it does not exist on directory backend
	_, err = suite.DirectoryBackEnd.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, backend.ErrNotFound)

	// Check it does not exist on underlayer backend
	_, err = suite.UnderlayerBackEnd.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, backend.ErrNotFound)
}

func (suite *DirectorySuite) TestRemoveDirectoryWhenDirectoryDoesNotExist() {
	// Remove the directory
	err := suite.Directory.RemoveDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, backend.ErrNotFound)
}

func (suite *DirectorySuite) TestRemoveDirectoryWhenOnlyOnUnderlayer() {
	// Create a directory on underlayer
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Remove the directory
	err = suite.Directory.RemoveDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Check it does not exist on directory backend
	_, err = suite.DirectoryBackEnd.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, backend.ErrNotFound)

	// Check it does not exist on underlayer backend
	_, err = suite.UnderlayerBackEnd.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, backend.ErrNotFound)
}
