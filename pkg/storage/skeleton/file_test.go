package skeleton

import (
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestFileSuite(t *testing.T) {
	t.Skip("Skeleton storage, not really implemented")
	suite.Run(t, new(FileSuite))
}

type FileSuite struct {
	test.FileSuite
}

func (s *FileSuite) SetupTest() {
	s.Directory = NewDirectory()
}
