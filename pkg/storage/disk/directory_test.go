package disk

import (
	"os"
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestDirectorySuite(t *testing.T) {
	suite.Run(t, new(DirectorySuite))
}

type DirectorySuite struct {
	path string
	test.DirectorySuite
}

func (suite *DirectorySuite) SetupTest() {
	// Create a temporary directory
	path, err := os.MkdirTemp("/tmp", "chonkfs-storage-disk-*")
	suite.Require().NoError(err)

	// Set it as the directory with underlayer as memory
	suite.Underlayer = mem.NewDirectory(nil)
	suite.Directory = NewDirectory(path, &DirectoryOptions{
		Underlayer: suite.Underlayer,
	})
	suite.path = path
}

func (suite *DirectorySuite) TearDownTest() {
	err := os.RemoveAll(suite.path)
	suite.Require().NoError(err)
}
