package mem

import (
	"testing"

	"github.com/lerenn/chonkfs/pkg/backends/test"
	"github.com/stretchr/testify/suite"
)

func TestFile(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

type FileSuite struct {
	test.FileSuite
}

func (suite *FileSuite) SetupTest() {
	suite.Directory = NewDirectory()
}
