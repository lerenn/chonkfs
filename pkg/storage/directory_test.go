package storage_test

import (
	"context"
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/lerenn/chonkfs/pkg/storage/backend"
	"github.com/lerenn/chonkfs/pkg/storage/backend/mem"
	"github.com/stretchr/testify/suite"
)

func TestDirectorySuite(t *testing.T) {
	suite.Run(t, new(DirectorySuite))
}

type DirectorySuite struct {
	DirectoryBackEnd  backend.BackEnd
	Directory         storage.Directory
	UnderlayerBackEnd backend.BackEnd
	Underlayer        storage.Directory
	suite.Suite
}

func (suite *DirectorySuite) SetupTest() {
	suite.UnderlayerBackEnd = mem.NewBackEnd()
	suite.Underlayer = storage.NewDirectory(suite.UnderlayerBackEnd, nil)

	suite.DirectoryBackEnd = mem.NewBackEnd()
	suite.Directory = storage.NewDirectory(suite.DirectoryBackEnd, &storage.DirectoryOptions{
		Underlayer: suite.Underlayer,
	})
}

func (suite *DirectorySuite) TestCreateDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Check it exists on directory backend
	err = suite.DirectoryBackEnd.IsDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Check it exists on underlayer backend
	err = suite.UnderlayerBackEnd.IsDirectory(context.Background(), "DirectoryA")
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

func (suite *DirectorySuite) TestInfo() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Get info
	info, err := suite.Directory.Info(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(storage.DirectoryInfo{}, info)
}

func (suite *DirectorySuite) TestListFiles() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateDirectory(context.Background(), "DirectoryB")
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateDirectory(context.Background(), "DirectoryC")
	suite.Require().NoError(err)

	// Create a file
	_, err = suite.Directory.CreateFile(context.Background(), "FileA", 4096)
	suite.Require().NoError(err)

	// List files
	files, err := suite.Directory.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(files, 3)
}
