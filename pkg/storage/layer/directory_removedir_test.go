package layer

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/storage"
)

func (suite *DirectorySuite) TestRemoveDirectoryOnBackendAndUnderlayer() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Remove the directory
	err = suite.Directory.RemoveDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Check it does not exist on directory backend
	_, err = suite.DirectoryBackEnd.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)

	// Check it does not exist on underlayer backend
	_, err = suite.UnderlayerBackEnd.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
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
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)

	// Check it does not exist on underlayer backend
	_, err = suite.UnderlayerBackEnd.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
}
