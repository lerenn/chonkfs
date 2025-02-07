package layer

import "context"

func (suite *DirectorySuite) TestGetDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Get the directory
	directory, err := suite.Directory.GetDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)
	suite.Require().NotNil(directory)
}

func (suite *DirectorySuite) TestGetDirectoryWhenOnlyOnUnderlayer() {
	// Create a directory on underlayer
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Get the directory
	directory, err := suite.Directory.GetDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)
	suite.Require().NotNil(directory)
}
