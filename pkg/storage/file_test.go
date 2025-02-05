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
	DirectoryBackEnd  backend.BackEnd
	Directory         storage.Directory
	UnderlayerBackEnd backend.BackEnd
	Underlayer        storage.Directory
	suite.Suite
}

func (suite *FileSuite) SetupTest() {
	suite.UnderlayerBackEnd = mem.NewBackEnd()
	suite.Underlayer = storage.NewDirectory(suite.UnderlayerBackEnd, nil)

	suite.DirectoryBackEnd = mem.NewBackEnd()
	suite.Directory = storage.NewDirectory(suite.DirectoryBackEnd, &storage.DirectoryOptions{
		Underlayer: suite.Underlayer,
	})
}

func (suite *FileSuite) TestCreateFile() {
	// Create a directory
	_, err := suite.Directory.CreateFile(context.Background(), "FileA", 4096)
	suite.Require().NoError(err)

	// Check it exists on directory backend
	err = suite.DirectoryBackEnd.IsFile(context.Background(), "FileA")
	suite.Require().NoError(err)

	// Check it exists on underlayer backend
	err = suite.UnderlayerBackEnd.IsFile(context.Background(), "FileA")
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
