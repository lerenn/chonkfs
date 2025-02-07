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
