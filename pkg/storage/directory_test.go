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
	DirectoryBackEnd  backend.Directory
	Directory         storage.Directory
	UnderlayerBackEnd backend.Directory
	Underlayer        storage.Directory
	suite.Suite
}

func (suite *DirectorySuite) SetupTest() {
	suite.UnderlayerBackEnd = mem.NewDirectory()
	suite.Underlayer = storage.NewDirectory(suite.UnderlayerBackEnd, nil)

	suite.DirectoryBackEnd = mem.NewDirectory()
	suite.Directory = storage.NewDirectory(suite.DirectoryBackEnd, &storage.DirectoryOptions{
		Underlayer: suite.Underlayer,
	})
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
