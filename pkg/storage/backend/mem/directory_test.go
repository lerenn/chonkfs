package mem

import (
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/backend/test"
	"github.com/stretchr/testify/suite"
)

func TestDirectorySuite(t *testing.T) {
	suite.Run(t, new(DirectorySuite))
}

type DirectorySuite struct {
	test.DirectorySuite
}

func (suite *DirectorySuite) SetupTest() {
	suite.Directory = NewDirectory()
}
