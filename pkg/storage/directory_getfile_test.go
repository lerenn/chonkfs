package storage_test

import "context"

func (suite *DirectorySuite) TestGetFile() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "File", 4096)
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "File")
	suite.Require().NoError(err)
	suite.Require().NotNil(file)
}

func (suite *DirectorySuite) TestGetFileWhenOnlyOnUnderlayer() {
	// Create a file on underlayer
	_, err := suite.Underlayer.CreateFile(context.Background(), "File", 4096)
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "File")
	suite.Require().NoError(err)
	suite.Require().NotNil(file)
}
