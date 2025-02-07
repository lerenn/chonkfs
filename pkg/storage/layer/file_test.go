package layer

import (
	"context"
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

type FileSuite struct {
	DirectoryBackEnd  storage.Directory
	UnderlayerBackEnd storage.Directory
	Underlayer        storage.Directory
	test.FileSuite
}

func (suite *FileSuite) SetupTest() {
	suite.UnderlayerBackEnd = mem.NewDirectory()
	suite.Underlayer = NewDirectory(suite.UnderlayerBackEnd, nil)

	suite.DirectoryBackEnd = mem.NewDirectory()
	suite.Directory = NewDirectory(suite.DirectoryBackEnd, &DirectoryOptions{
		Underlayer: suite.Underlayer,
	})
}

func (suite *FileSuite) TestCreateFileAndCheckUnderlayer() {
	// Create a directory
	_, err := suite.Directory.CreateFile(context.Background(), "FileA", 4096)
	suite.Require().NoError(err)

	// Check it exists on underlayer backend
	_, err = suite.UnderlayerBackEnd.GetFile(context.Background(), "FileA")
	suite.Require().NoError(err)
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
