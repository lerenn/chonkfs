package skeleton

import (
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestDirectorySuite(t *testing.T) {
	t.Skip("skipping test as it is a skeleton")
	suite.Run(t, new(DirectorySuite))
}

type DirectorySuite struct {
	test.DirectorySuite
}

func (suite *DirectorySuite) SetupTest() {
	suite.Directory = NewDirectory()
}
