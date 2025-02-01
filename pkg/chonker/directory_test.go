package chonker

import (
	"context"
	"testing"

	"github.com/lerenn/chonkfs/pkg/storage/mem"
	"github.com/stretchr/testify/suite"
)

func TestDirectorySuite(t *testing.T) {
	suite.Run(t, new(DirectorySuite))
}

type DirectorySuite struct {
	suite.Suite
}

func (suite *DirectorySuite) TestListFiles() {
	// Create a directory
	d, err := NewDirectory(context.Background(), mem.NewDirectory(nil))
	suite.Require().NoError(err)

	// Create a file
	_, err = d.CreateFile(context.Background(), "FileA.txt", 4)
	suite.Require().NoError(err)

	// Create a file
	_, err = d.CreateFile(context.Background(), "FileB.txt", 4)
	suite.Require().NoError(err)

	// List files
	files, err := d.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(files, 2)
}
