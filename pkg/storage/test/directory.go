package test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/stretchr/testify/suite"
)

type DirectorySuite struct {
	Directory storage.Directory
	suite.Suite
}

func (suite *DirectorySuite) TestCreateGetDirectory() {
	// Create a directory
	nd, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Read the directory
	dir, err := suite.Directory.GetDirectory(context.Background(), "dir")
	suite.Require().NoError(err)
	suite.Require().Equal(nd, dir)
}

func (suite *DirectorySuite) TestListDirectories() {
	// Create three directories
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir1")
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateDirectory(context.Background(), "dir2")
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateDirectory(context.Background(), "dir3")
	suite.Require().NoError(err)

	// List directories
	dirs, err := suite.Directory.ListDirectories(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(dirs, 3)
	suite.Require().Contains(dirs, "dir1")
	suite.Require().Contains(dirs, "dir2")
	suite.Require().Contains(dirs, "dir3")
}

func (suite *DirectorySuite) TestInfo() {
	info, err := suite.Directory.Info(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(storage.DirectoryInfo{}, info)
}

func (suite *DirectorySuite) TestRemoveDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Remove the directory
	err = suite.Directory.RemoveDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Check if the directory is removed
	_, err = suite.Directory.GetDirectory(context.Background(), "dir")
	suite.Require().Error(err)

	// Check if the directory is not in list
	dirs, err := suite.Directory.ListDirectories(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(dirs, 0)
}

func (suite *DirectorySuite) TestRenameDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Rename the directory
	err = suite.Directory.RenameDirectory(context.Background(), "dir", suite.Directory, "newdir", false)
	suite.Require().NoError(err)

	// Check if the directory is renamed
	dir, err := suite.Directory.GetDirectory(context.Background(), "newdir")
	suite.Require().NoError(err)
	suite.Require().NotNil(dir)

	// Check if the old directory is removed
	_, err = suite.Directory.GetDirectory(context.Background(), "dir")
	suite.Require().Error(err)
}

func (suite *DirectorySuite) TestCreateGetFile() {
	// Create a file
	f, err := suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().NoError(err)

	// Read the file
	file, err := suite.Directory.GetFile(context.Background(), "file")
	suite.Require().NoError(err)
	suite.Require().Equal(f, file)
}

func (suite *DirectorySuite) TestListFiles() {
	// Create three files
	_, err := suite.Directory.CreateFile(context.Background(), "file1", 1)
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "file2", 1)
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "file3", 1)
	suite.Require().NoError(err)

	// List files
	files, err := suite.Directory.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(files, 3)
	suite.Require().Contains(files, "file1")
	suite.Require().Contains(files, "file2")
	suite.Require().Contains(files, "file3")
}

func (suite *DirectorySuite) TestCreateFileAlreadyExists() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().NoError(err)

	// Create the same file
	_, err = suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().Error(err)
}

func (suite *DirectorySuite) TestRemoveFile() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().NoError(err)

	// Remove the file
	err = suite.Directory.RemoveFile(context.Background(), "file")
	suite.Require().NoError(err)

	// Check if the file is removed
	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().Error(err)

	// Check if the file is not in list
	files, err := suite.Directory.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(files, 0)
}

func (suite *DirectorySuite) TestRenameFile() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	// Rename the file
	err = suite.Directory.RenameFile(context.Background(), "file", suite.Directory, "newfile", false)
	suite.Require().NoError(err)

	// Check if the file is renamed
	file, err := suite.Directory.GetFile(context.Background(), "newfile")
	suite.Require().NoError(err)
	suite.Require().NotNil(file)

	// Check if the old file is removed
	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().Error(err)
}
