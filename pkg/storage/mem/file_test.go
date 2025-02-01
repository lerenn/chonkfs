package mem

import (
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/test"
	"github.com/stretchr/testify/suite"
)

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

type FileSuite struct {
	test.FileSuite
}

func (s *FileSuite) SetupTest() {
	s.Underlayer = NewDirectory(nil)
	s.Directory = NewDirectory(&DirectoryOptions{
		Underlayer: s.Underlayer,
	})
}
