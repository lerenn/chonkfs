package layer

import (
	"context"
	"testing"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestDirectorySuite(t *testing.T) {
	suite.Run(t, new(DirectorySuite))
}

type DirectorySuite struct {
	DirectoryBackEnd  storage.Directory
	UnderlayerBackEnd storage.Directory
	Underlayer        storage.Directory
	test.DirectorySuite
}

func (suite *DirectorySuite) SetupTest() {
	suite.UnderlayerBackEnd = mem.NewDirectory()
	suite.Underlayer = NewDirectory(suite.UnderlayerBackEnd, nil)

	suite.DirectoryBackEnd = mem.NewDirectory()
	suite.Directory = NewDirectory(suite.DirectoryBackEnd, &DirectoryOptions{
		Underlayer: suite.Underlayer,
	})
}

func (suite *DirectorySuite) TestGetInfoWhenDirectoryExistsOnlyOnUnderlayer() {
	// Create a directory
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Get directory
	dir, err := suite.Directory.GetDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Get info
	dirInfo, err := dir.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(info.Directory{}, dirInfo)
}

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

func (suite *DirectorySuite) TestListFilesWithOneInUnderlayer() {
	// Create a directory
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Create a file in underlayer
	_, err = suite.Underlayer.CreateFile(context.Background(), "FileA", 4096)
	suite.Require().NoError(err)

	// Create 2 files in directory
	_, err = suite.Directory.CreateFile(context.Background(), "FileB", 4096)
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "FileC", 4096)
	suite.Require().NoError(err)

	// List files
	files, err := suite.Directory.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(files, 3)
}

func (suite *DirectorySuite) TestGetFileWhenOnlyOnUnderlayer() {
	// Create a file on underlayer
	_, err := suite.Underlayer.CreateFile(context.Background(), "File", 4096)
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "File")
	suite.Require().NoError(err)
	suite.Require().NotNil(file)
}

func (suite *DirectorySuite) TestGetDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Get the directory
	directory, err := suite.Directory.GetDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)
	suite.Require().NotNil(directory)
}

func (suite *DirectorySuite) TestGetDirectoryWhenOnlyOnUnderlayer() {
	// Create a directory on underlayer
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Get the directory
	directory, err := suite.Directory.GetDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)
	suite.Require().NotNil(directory)
}

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
