package skeleton

import (
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestDirectorySuite(t *testing.T) {
	t.Skip("Skeleton storage, not really implemented")
	suite.Run(t, new(DirectorySuite))
}

type DirectorySuite struct {
	test.DirectorySuite
}

func (suite *DirectorySuite) SetupTest() {
	suite.Underlayer = mem.NewDirectory(nil)
	suite.Directory = NewDirectory(&DirectoryOptions{
		Underlayer: suite.Underlayer,
	})
}
