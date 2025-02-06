package storage_test

import (
	"context"
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/lerenn/chonkfs/pkg/storage/backend"
	"github.com/lerenn/chonkfs/pkg/storage/backend/mem"
	"github.com/stretchr/testify/suite"
)

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

type FileSuite struct {
	DirectoryBackEnd  backend.Directory
	Directory         storage.Directory
	UnderlayerBackEnd backend.Directory
	Underlayer        storage.Directory
	suite.Suite
}

func (suite *FileSuite) SetupTest() {
	suite.UnderlayerBackEnd = mem.NewDirectory()
	suite.Underlayer = storage.NewDirectory(suite.UnderlayerBackEnd, nil)

	suite.DirectoryBackEnd = mem.NewDirectory()
	suite.Directory = storage.NewDirectory(suite.DirectoryBackEnd, &storage.DirectoryOptions{
		Underlayer: suite.Underlayer,
	})
}

func (suite *FileSuite) TestCreateFile() {
	// Create a directory
	_, err := suite.Directory.CreateFile(context.Background(), "FileA", 4096)
	suite.Require().NoError(err)

	// Check it exists on directory backend
	_, err = suite.DirectoryBackEnd.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)

	// Check it exists on underlayer backend
	_, err = suite.UnderlayerBackEnd.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
}

func (suite *FileSuite) TestCreateFileWhenFileAlreadyExists() {
	// Create a directory
	_, err := suite.Directory.CreateFile(context.Background(), "FileA", 4096)
	suite.Require().NoError(err)

	// Create the same directory again
	_, err = suite.Directory.CreateFile(context.Background(), "FileA", 4096)
	suite.Require().ErrorIs(err, storage.ErrFileAlreadyExists)
}

func (suite *FileSuite) TestCreateFileWhenDirectoryAlreadyExists() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Create a file with the same name
	_, err = suite.Directory.CreateFile(context.Background(), "DirectoryA", 4096)
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}

func (suite *FileSuite) TestGetInfo() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "FileA", 4096)
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	suite.Require().NotNil(file)

	// Get the file info
	info, err := file.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(4096, info.ChunkSize)
}

func (suite *FileSuite) TestGetInfoWhenFileExistsOnlyOnUnderlayer() {
	// Create a file on underlayer
	_, err := suite.Underlayer.CreateFile(context.Background(), "FileA", 4096)
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
	suite.Require().NotNil(file)

	// Get the file info
	info, err := file.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(4096, info.ChunkSize)
}
