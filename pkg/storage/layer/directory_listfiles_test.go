package layer

import (
	"context"
)

func (suite *DirectorySuite) TestListFilesWithOneInUnderlayer() {
	// Create a directory
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Create a file in underlayer
	_, err = suite.Underlayer.CreateFile(context.Background(), "FileA", 4096)
	suite.Require().NoError(err)

	// Create 2 files in directory
	_, err = suite.Directory.CreateFile(context.Background(), "FileB", 4096)
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "FileC", 4096)
	suite.Require().NoError(err)

	// List files
	files, err := suite.Directory.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(files, 3)
}
