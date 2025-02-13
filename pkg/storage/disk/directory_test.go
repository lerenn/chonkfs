package disk

import (
	"os"
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestDirectorySuite(t *testing.T) {
	suite.Run(t, new(DirectorySuite))
}

type DirectorySuite struct {
	test.DirectorySuite
	Path string
}

func (suite *DirectorySuite) SetupTest() {
	path, err := os.MkdirTemp("", "chonkfs-test-*")
	suite.Require().NoError(err)
	suite.Path = path
	suite.Directory = NewDirectory(path)
}

func (suite *DirectorySuite) TearDownTest() {
	err := os.RemoveAll(suite.Path)
	suite.Require().NoError(err)
}
