package test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/stretchr/testify/suite"
)

// DirectorySuite is a test suite for a directory.
type DirectorySuite struct {
	Underlayer storage.Directory
	Directory  storage.Directory
	suite.Suite
}

// TestGetInexistantDirectory checks getting an inexistant directory.
func (suite *DirectorySuite) TestGetInexistantDirectory() {
	// Read the directory
	_, err := suite.Directory.GetDirectory(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotExists)
}

// TestCreateGetDirectory tests creating and getting a directory.
func (suite *DirectorySuite) TestCreateGetDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Read the directory
	_, err = suite.Directory.GetDirectory(context.Background(), "dir")
	suite.Require().NoError(err)
}

// TestCreateGetUnderlayerDirectory tests creating and getting a directory in the underlayer.
func (suite *DirectorySuite) TestCreateGetUnderlayerDirectory() {
	// Create a directory
	d, err := suite.Directory.CreateDirectory(context.Background(), "dir1")
	suite.Require().NoError(err)

	// Read the directory in the underlayer
	ud, err := suite.Underlayer.GetDirectory(context.Background(), "dir1")
	suite.Require().NoError(err)

	// Create a directory in the child directory
	_, err = d.CreateDirectory(context.Background(), "dir2")
	suite.Require().NoError(err)

	// Read the child's child directory in the underlayer
	_, err = ud.GetDirectory(context.Background(), "dir2")
	suite.Require().NoError(err)
}

// TestGetThenCreateUnderlayerDirectory tests if when getting a directory, writing
// it will write the underlayer directory as well.
func (suite *DirectorySuite) TestGetThenCreateUnderlayerDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir1")
	suite.Require().NoError(err)

	// Read the directory in the underlayer
	ud, err := suite.Underlayer.GetDirectory(context.Background(), "dir1")
	suite.Require().NoError(err)

	// Read the directory
	d, err := suite.Directory.GetDirectory(context.Background(), "dir1")
	suite.Require().NoError(err)

	// Create a directory in the child directory
	_, err = d.CreateDirectory(context.Background(), "dir2")
	suite.Require().NoError(err)

	// Read the child's child directory in the underlayer
	_, err = ud.GetDirectory(context.Background(), "dir2")
	suite.Require().NoError(err)
}

// TestUnderlayerDirectoryExistsAndFillAbove checks if when a directory exists in the
// underlayer then the directory exists above as well.
func (suite *DirectorySuite) TestUnderlayerDirectoryExistsAndFillAbove() {
	// Create the directory in the underlayer
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Read it from above
	_, err = suite.Directory.GetDirectory(context.Background(), "dir")
	suite.Require().NoError(err)
}

// TestListDirectories tests listing directories.
func (suite *DirectorySuite) TestListDirectories() {
	// Create three directories
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir1")
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateDirectory(context.Background(), "dir2")
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateDirectory(context.Background(), "dir3")
	suite.Require().NoError(err)

	// Create a file
	_, err = suite.Directory.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	// List directories
	dirs, err := suite.Directory.ListDirectories(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(dirs, 3)
	suite.Require().Contains(dirs, "dir1")
	suite.Require().Contains(dirs, "dir2")
	suite.Require().Contains(dirs, "dir3")
}

// TestListDirectoriesThatExistsOnUnderlayerOnly checks that when a directory
// exists on the underlayer, then it will be correctly displayed on the listing
// above.
func (suite *DirectorySuite) TestListDirectoriesThatExistsOnUnderlayerOnly() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir_above")
	suite.Require().NoError(err)

	// Create a directory on the underlayer only
	_, err = suite.Underlayer.CreateDirectory(context.Background(), "underlayer")
	suite.Require().NoError(err)

	// List directories
	dirs, err := suite.Directory.ListDirectories(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(dirs, 2)
	suite.Require().Contains(dirs, "dir_above")
	suite.Require().Contains(dirs, "underlayer")
}

// TestInfo tests getting the directory info.
func (suite *DirectorySuite) TestInfo() {
	info, err := suite.Directory.Info(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(storage.DirectoryInfo{}, info)
}

// TestCreateDirectoryAlreadyExists tests creating a directory that already exists.
func (suite *DirectorySuite) TestCreateDirectoryAlreadyExists() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Create the same directory
	_, err = suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}

// TestRemoveDirectory tests removing a directory.
func (suite *DirectorySuite) TestRemoveDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Remove the directory
	err = suite.Directory.RemoveDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Check if the directory is removed
	_, err = suite.Directory.GetDirectory(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotExists)

	// Check if the directory is not in list
	dirs, err := suite.Directory.ListDirectories(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(dirs, 0)
}

func (suite *DirectorySuite) RemoveInexistantDirectory() {
	err := suite.Directory.RemoveDirectory(context.Background(), "inexistant")
	suite.Require().NoError(err)
}

// TestRemoveDirectoryAndCheckUnderlayer tests removing a directory is passed to underlayer.
func (suite *DirectorySuite) TestRemoveDirectoryAndCheckUnderlayer() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Remove the directory
	err = suite.Directory.RemoveDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Check if the directory is removed in the underlayer
	_, err = suite.Underlayer.GetDirectory(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotExists)

	// Check if the directory is not in list in the underlayer
	udirs, err := suite.Underlayer.ListDirectories(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(udirs, 0)
}

// TestRenameDirectory tests renaming a directory.
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
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotExists)
}

// TestRenameDirectoryAndCheckUnderlayer tests renaming a directory is passed to underlayer.
func (suite *DirectorySuite) TestRenameDirectoryAndCheckUnderlayer() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Rename the directory
	err = suite.Directory.RenameDirectory(context.Background(), "dir", suite.Directory, "newdir", false)
	suite.Require().NoError(err)

	// Check if the directory is renamed in the underlayer
	udir, err := suite.Underlayer.GetDirectory(context.Background(), "newdir")
	suite.Require().NoError(err)
	suite.Require().NotNil(udir)

	// Check if the old directory is removed in the underlayer
	_, err = suite.Underlayer.GetDirectory(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotExists)
}

// TestRenameDirectoryThatExistOnlyOnUnderlayer checks that when a directory
// that only exists on the underlayer, then it is renamed and present above
func (suite *DirectorySuite) TestRenameDirectoryThatExistOnlyOnUnderlayer() {
	// Create a directory on the underlayer
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "dir")
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
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotExists)
}

// TestCreateGetFile tests creating and getting a file.
func (suite *DirectorySuite) TestCreateGetFile() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().NoError(err)

	// Read the file
	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().NoError(err)
}

// TestCreateGetFileUnderlayer tests creating and getting a file from underlayer.
func (suite *DirectorySuite) TestCreateGetFileUnderlayer() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().NoError(err)

	// Read the file in the underlayer
	_, err = suite.Underlayer.GetFile(context.Background(), "file")
	suite.Require().NoError(err)
}

// TestUnderlayerFileExistsAndFillAbove checks if when a directory exists in the
// underlayer then the file exists above as well.
func (suite *DirectorySuite) TestUnderlayerFileExistsAndFillAbove() {
	// Create the file in the underlayer
	_, err := suite.Underlayer.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	// Read it from above
	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().NoError(err)
}

// TestListFiles tests listing files.
func (suite *DirectorySuite) TestListFiles() {
	// Create three files
	_, err := suite.Directory.CreateFile(context.Background(), "file1", 1)
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "file2", 1)
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "file3", 1)
	suite.Require().NoError(err)

	// Create one directory
	_, err = suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// List files
	files, err := suite.Directory.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(files, 3)
	suite.Require().Contains(files, "file1")
	suite.Require().Contains(files, "file2")
	suite.Require().Contains(files, "file3")
}

// TestListFilesThatExistsOnUnderlayerOnly checks that when a file
// exists on the underlayer, then it will be correctly displayed on the listing
// above.
func (suite *DirectorySuite) TestListFilesThatExistsOnUnderlayerOnly() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "file_above", 4096)
	suite.Require().NoError(err)

	// Create a file on the underlayer only
	_, err = suite.Underlayer.CreateFile(context.Background(), "underlayer", 4096)
	suite.Require().NoError(err)

	// List files
	dirs, err := suite.Directory.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(dirs, 2)
	suite.Require().Contains(dirs, "file_above")
	suite.Require().Contains(dirs, "underlayer")
}

// TestCreateFileAlreadyExists tests creating a file that already exists.
func (suite *DirectorySuite) TestCreateFileAlreadyExists() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().NoError(err)

	// Create the same file
	_, err = suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().ErrorIs(err, storage.ErrFileAlreadyExists)
}

// TestRemoveFile tests removing a file.
func (suite *DirectorySuite) TestRemoveFile() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().NoError(err)

	// Remove the file
	err = suite.Directory.RemoveFile(context.Background(), "file")
	suite.Require().NoError(err)

	// Check if the file is removed
	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().ErrorIs(err, storage.ErrFileNotExists)

	// Check if the file is not in list
	files, err := suite.Directory.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(files, 0)
}

// TestRemoveFileAndCheckUnderlayer tests removing a file is passed ti underlayer.
func (suite *DirectorySuite) TestRemoveFileAndCheckUnderlayer() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "file", 1)
	suite.Require().NoError(err)

	// Remove the file
	err = suite.Directory.RemoveFile(context.Background(), "file")
	suite.Require().NoError(err)

	// Check if the file is removed in the underlayer
	_, err = suite.Underlayer.GetFile(context.Background(), "file")
	suite.Require().ErrorIs(err, storage.ErrFileNotExists)

	// Check if the file is not in list in the underlayer
	ufiles, err := suite.Underlayer.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(ufiles, 0)
}

// // TestRenameFile tests renaming a file.
// func (suite *DirectorySuite) TestRenameFile() {
// 	// Create a file
// 	_, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
// 	suite.Require().NoError(err)

// 	// Rename the file
// 	err = suite.Directory.RenameFile(context.Background(), "file", suite.Directory, "newfile", false)
// 	suite.Require().NoError(err)

// 	// Check if the file is renamed
// 	file, err := suite.Directory.GetFile(context.Background(), "newfile")
// 	suite.Require().NoError(err)
// 	suite.Require().NotNil(file)

// 	// Check if the old file is removed
// 	_, err = suite.Directory.GetFile(context.Background(), "file")
// 	suite.Require().ErrorIs(err, storage.ErrFileNotExists)
// }

// // TestRenameFileAndCheckUnderlayer tests renaming a file is passed to underlayer.
// func (suite *DirectorySuite) TestRenameFileAndCheckUnderlayer() {
// 	// Create a file
// 	_, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
// 	suite.Require().NoError(err)

// 	// Rename the file
// 	err = suite.Directory.RenameFile(context.Background(), "file", suite.Directory, "newfile", false)
// 	suite.Require().NoError(err)

// 	// Check if the file is renamed in the underlayer
// 	ufile, err := suite.Underlayer.GetFile(context.Background(), "newfile")
// 	suite.Require().NoError(err)
// 	suite.Require().NotNil(ufile)

// 	// Check if the old file is removed in the underlayer
// 	_, err = suite.Underlayer.GetFile(context.Background(), "file")
// 	suite.Require().ErrorIs(err, storage.ErrFileNotExists)
// }

// TestRenameFileThatExistOnlyOnUnderlayer checks that when a directory
// that only exists on the underlayer, then it is renamed and present above
// func (suite *DirectorySuite) TestRenameFileThatExistOnlyOnUnderlayer() {
// 	// Create a directory on the underlayer
// 	_, err := suite.Underlayer.CreateFile(context.Background(), "file", 4096)
// 	suite.Require().NoError(err)

// 	// Rename the directory
// 	err = suite.Directory.RenameFile(context.Background(), "file", suite.Directory, "newfile", false)
// 	suite.Require().NoError(err)

// 	// Check if the directory is renamed
// 	dir, err := suite.Directory.GetFile(context.Background(), "newfile")
// 	suite.Require().NoError(err)
// 	suite.Require().NotNil(dir)

// 	// Check if the old directory is removed
// 	_, err = suite.Directory.GetFile(context.Background(), "file")
// 	suite.Require().ErrorIs(err, storage.ErrDirectoryNotExists)
// }
