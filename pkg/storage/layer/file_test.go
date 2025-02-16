package layer

import (
	"os"
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/disk"
	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/stretchr/testify/suite"
)

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileWithMemSuite))
	suite.Run(t, new(FileWithDiskSuite))
}

type FileWithMemSuite struct {
	FileSuite
}

func (suite *FileWithMemSuite) SetupTest() {
	suite.Upperlayer = mem.NewDirectory()
	suite.Underlayer = mem.NewDirectory()
	suite.Directory, _ = NewDirectory(suite.Upperlayer, suite.Underlayer)
}

type FileWithDiskSuite struct {
	UpperlayerPath string
	UnderlayerPath string
	FileSuite
}

func (suite *FileWithDiskSuite) SetupTest() {
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

func (suite *FileWithDiskSuite) TearDownTest() {
	err := os.RemoveAll(suite.UpperlayerPath)
	suite.Require().NoError(err)
	err = os.RemoveAll(suite.UnderlayerPath)
	suite.Require().NoError(err)
}
