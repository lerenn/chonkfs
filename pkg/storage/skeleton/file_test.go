package skeleton

import (
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestFileSuite(t *testing.T) {
	t.Skip("skipping test as it is a skeleton")
	suite.Run(t, new(FileSuite))
}

type FileSuite struct {
	test.FileSuite
}

func (suite *FileSuite) SetupTest() {
	suite.Directory = NewDirectory()
}
