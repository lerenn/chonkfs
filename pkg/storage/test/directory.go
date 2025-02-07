package test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/stretchr/testify/suite"
)

type DirectorySuite struct {
	Directory storage.Directory
	suite.Suite
}

func (suite *DirectorySuite) TestCreateDirectory() {
	d, err := suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().NoError(err)
	suite.Require().NotNil(d)

	rd, err := suite.Directory.GetDirectory(context.Background(), "toto")
	suite.Require().NoError(err)
	suite.Require().NotNil(rd)
}

func (suite *DirectorySuite) TestCreateDirectoryOnExistingFile() {
	_, err := suite.Directory.CreateFile(context.Background(), "toto", 4096)
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, storage.ErrFileAlreadyExists)
}

func (suite *DirectorySuite) TestCreateDirectoryOnExistingDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}

func (suite *DirectorySuite) TestGetDirectoryWhenDoesNotExist() {
	_, err := suite.Directory.GetDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
}

func (suite *DirectorySuite) TestGetDirectoryWhenIsFile() {
	_, err := suite.Directory.CreateFile(context.Background(), "toto", 4096)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, storage.ErrIsFile)
}

func (suite *DirectorySuite) TestListFiles() {
	_, err := suite.Directory.CreateFile(context.Background(), "1", 4096)
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "2", 4096)
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "3", 4096)
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	// Check content
	files, err := suite.Directory.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Contains(files, "1")
	suite.Require().Contains(files, "2")
	suite.Require().Contains(files, "3")

	// Check length
	infoFile1, err := files["1"].GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(4096, infoFile1.ChunkSize)
}

func (suite *DirectorySuite) TestRemoveDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	err = suite.Directory.RemoveDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	_, err = suite.Directory.GetDirectory(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
}

func (suite *DirectorySuite) TestRemoveDirectoryWhenDoesNotExist() {
	err := suite.Directory.RemoveDirectory(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
}

func (suite *DirectorySuite) TestGetInfo() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Get info
	dirInfo, err := suite.Directory.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(info.Directory{}, dirInfo)
}

func (suite *DirectorySuite) TestGetFile() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "File", 4096)
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "File")
	suite.Require().NoError(err)
	suite.Require().NotNil(file)
}

func (suite *DirectorySuite) TestListDirectories() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Create a directory
	_, err = suite.Directory.CreateDirectory(context.Background(), "DirectoryB")
	suite.Require().NoError(err)

	// List directories
	dirs, err := suite.Directory.ListDirectories(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(dirs, 2)
}

func (suite *DirectorySuite) TestRemoveFile() {
	_, err := suite.Directory.CreateFile(context.Background(), "dir", 4096)
	suite.Require().NoError(err)

	err = suite.Directory.RemoveFile(context.Background(), "dir")
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)
}

func (suite *DirectorySuite) TestRemoveFileWhenDoesNotExist() {
	err := suite.Directory.RemoveFile(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)
}

func (suite *DirectorySuite) TestRenameFileOnSameDirectory() {
	_, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	err = suite.Directory.RenameFile(context.Background(), "file", suite.Directory, "newFile", true)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	_, err = suite.Directory.GetFile(context.Background(), "newFile")
	suite.Require().NoError(err)
}

func (suite *DirectorySuite) TestRenameFileOnDifferentDirectory() {
	_, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	dir, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	err = suite.Directory.RenameFile(context.Background(), "file", dir, "newFile", true)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	_, err = dir.GetFile(context.Background(), "newFile")
	suite.Require().NoError(err)
}

func (suite *DirectorySuite) TestRenameFileOnExistingFileWithNoReplace() {
	_, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateFile(context.Background(), "newFile", 4096)
	suite.Require().NoError(err)

	err = suite.Directory.RenameFile(context.Background(), "file", suite.Directory, "newFile", false)
	suite.Require().ErrorIs(err, storage.ErrFileAlreadyExists)
}

func (suite *DirectorySuite) TestRenameFileOnExistingFileWithReplace() {
	_, err := suite.Directory.CreateFile(context.Background(), "file", 4096)
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateFile(context.Background(), "newFile", 8192)
	suite.Require().NoError(err)

	err = suite.Directory.RenameFile(context.Background(), "file", suite.Directory, "newFile", true)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	f, err := suite.Directory.GetFile(context.Background(), "newFile")
	suite.Require().NoError(err)

	info, err := f.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(4096, info.ChunkSize)
}

func (suite *DirectorySuite) TestRenameDirectoryOnSameDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "directory")
	suite.Require().NoError(err)

	err = suite.Directory.RenameDirectory(context.Background(), "directory", suite.Directory, "newDirectory", true)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetDirectory(context.Background(), "directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)

	_, err = suite.Directory.GetDirectory(context.Background(), "newDirectory")
	suite.Require().NoError(err)
}

func (suite *DirectorySuite) TestRenameDirectoryOnDifferentDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "directory")
	suite.Require().NoError(err)

	dir, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	err = suite.Directory.RenameDirectory(context.Background(), "directory", dir, "newDirectory", true)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetDirectory(context.Background(), "directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)

	_, err = dir.GetDirectory(context.Background(), "newDirectory")
	suite.Require().NoError(err)
}

func (suite *DirectorySuite) TestRenameDirectoryOnExistingDirectoryWithNoReplace() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "directory")
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "newDirectory")
	suite.Require().NoError(err)

	err = suite.Directory.RenameDirectory(context.Background(), "directory", suite.Directory, "newDirectory", false)
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}

func (suite *DirectorySuite) TestRenameDirectoryOnExistingDirectoryWithReplace() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "directory")
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "newDirectory")
	suite.Require().NoError(err)

	err = suite.Directory.RenameDirectory(context.Background(), "directory", suite.Directory, "newDirectory", true)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetDirectory(context.Background(), "directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
}
