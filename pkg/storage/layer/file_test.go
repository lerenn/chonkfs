package layer

import (
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/stretchr/testify/suite"
)

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileWithMemSuite))
}

type FileWithMemSuite struct {
	FileSuite
}

func (suite *FileWithMemSuite) SetupTest() {
	suite.Upperlayer = mem.NewDirectory()
	suite.Underlayer = mem.NewDirectory()
	suite.Directory = NewDirectory(suite.Upperlayer, suite.Underlayer)
}
