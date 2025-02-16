package layer

import (
	"os"
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/disk"
	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/stretchr/testify/suite"
)

func TestDirectorySuite(t *testing.T) {
	suite.Run(t, new(DirectoryWithMemSuite))
	suite.Run(t, new(DirectoryWithDiskSuite))
}

type DirectoryWithMemSuite struct {
	DirectorySuite
}

func (suite *DirectoryWithMemSuite) SetupTest() {
	suite.Upperlayer = mem.NewDirectory()
	suite.Underlayer = mem.NewDirectory()
	suite.Directory, _ = NewDirectory(suite.Upperlayer, suite.Underlayer)
}

type DirectoryWithDiskSuite struct {
	UpperlayerPath string
	UnderlayerPath string
	DirectorySuite
}

func (suite *DirectoryWithDiskSuite) SetupTest() {
	path, err := os.MkdirTemp("", "chonkfs-underlayer-test-*")
	suite.Require().NoError(err)
	suite.UnderlayerPath = path

	path, err = os.MkdirTemp("", "chonkfs-upperlayer-test-*")
	suite.Require().NoError(err)
	suite.UpperlayerPath = path

	suite.Upperlayer = disk.NewDirectory(suite.UpperlayerPath)
	suite.Underlayer = disk.NewDirectory(suite.UnderlayerPath)
	suite.Directory, _ = NewDirectory(suite.Upperlayer, suite.Underlayer)
}

func (suite *DirectoryWithDiskSuite) TearDownTest() {
	err := os.RemoveAll(suite.UpperlayerPath)
	suite.Require().NoError(err)
	err = os.RemoveAll(suite.UnderlayerPath)
	suite.Require().NoError(err)
}
