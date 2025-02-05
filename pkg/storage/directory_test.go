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

func (suite *DirectorySuite) TestCreateDirectoryAlreadyExists() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Create the same directory again
	_, err = suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}
