package disk

import (
	"os"
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

type FileSuite struct {
	test.FileSuite
	Path string
}

func (suite *FileSuite) SetupTest() {
	path, err := os.MkdirTemp("", "chonkfs-test-*")
	suite.Require().NoError(err)
	suite.Path = path
	suite.Directory = NewDirectory(path)
}

func (suite *FileSuite) TearDownTest() {
	err := os.RemoveAll(suite.Path)
	suite.Require().NoError(err)
}
