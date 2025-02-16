package test

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/stretchr/testify/suite"
)

// DirectorySuite is a test suite for storage.Directory implementations.
type DirectorySuite struct {
	Directory storage.Directory
	suite.Suite
}

// TestCreateDirectory tests the creation of a directory.
func (suite *DirectorySuite) TestCreateDirectory() {
	d, err := suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().NoError(err)
	suite.Require().NotNil(d)

	rd, err := suite.Directory.GetDirectory(context.Background(), "toto")
	suite.Require().NoError(err)
	suite.Require().NotNil(rd)
}

// TestCreateDirectoryOnExistingFile tests the creation of a directory on an existing file.
func (suite *DirectorySuite) TestCreateDirectoryOnExistingFile() {
	_, err := suite.Directory.CreateFile(context.Background(), "toto", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, storage.ErrFileAlreadyExists)
}

// TestCreateDirectoryOnExistingDirectory tests the creation of a directory on an existing directory.
func (suite *DirectorySuite) TestCreateDirectoryOnExistingDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}

// TestGetDirectoryWhenDoesNotExist tests the retrieval of a directory that does not exist.
func (suite *DirectorySuite) TestGetDirectoryWhenDoesNotExist() {
	_, err := suite.Directory.GetDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
}

// TestGetDirectoryWhenIsFile tests the retrieval of a directory that is a file.
func (suite *DirectorySuite) TestGetDirectoryWhenIsFile() {
	_, err := suite.Directory.CreateFile(context.Background(), "toto", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	_, err = suite.Directory.GetDirectory(context.Background(), "toto")
	suite.Require().ErrorIs(err, storage.ErrIsFile)
}

// TestListFiles tests the listing of files.
func (suite *DirectorySuite) TestListFiles() {
	_, err := suite.Directory.CreateFile(context.Background(), "1", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "2", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "3", info.File{
		ChunkSize: 4096,
	})
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

// TestRemoveDirectory tests the removal of a directory.
func (suite *DirectorySuite) TestRemoveDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	err = suite.Directory.RemoveDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	_, err = suite.Directory.GetDirectory(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
}

// TestRemoveDirectoryWhenDoesNotExist tests the removal of a directory that does not exist.
func (suite *DirectorySuite) TestRemoveDirectoryWhenDoesNotExist() {
	err := suite.Directory.RemoveDirectory(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
}

// TestGetInfo tests the retrieval of directory information.
func (suite *DirectorySuite) TestGetInfo() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Get info
	dirInfo, err := suite.Directory.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(info.Directory{}, dirInfo)
}

// TestGetFile tests the retrieval of a file.
func (suite *DirectorySuite) TestGetFile() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "File", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "File")
	suite.Require().NoError(err)
	suite.Require().NotNil(file)
}

// TestListDirectories tests the listing of directories.
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

// TestRemoveFile tests the removal of a file.
func (suite *DirectorySuite) TestRemoveFile() {
	_, err := suite.Directory.CreateFile(context.Background(), "dir", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	err = suite.Directory.RemoveFile(context.Background(), "dir")
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)
}

// TestRemoveFileWhenDoesNotExist tests the removal of a file that does not exist.
func (suite *DirectorySuite) TestRemoveFileWhenDoesNotExist() {
	err := suite.Directory.RemoveFile(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)
}

// TestRenameFileOnSameDirectory tests the renaming of a file on the same directory.
func (suite *DirectorySuite) TestRenameFileOnSameDirectory() {
	_, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	err = suite.Directory.RenameFile(context.Background(), "file", suite.Directory, "newFile", false)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	_, err = suite.Directory.GetFile(context.Background(), "newFile")
	suite.Require().NoError(err)
}

// TestRenameFileOnDifferentDirectory tests the renaming of a file on a different directory.
func (suite *DirectorySuite) TestRenameFileOnDifferentDirectory() {
	_, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	dir, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	err = suite.Directory.RenameFile(context.Background(), "file", dir, "newFile", false)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	_, err = dir.GetFile(context.Background(), "newFile")
	suite.Require().NoError(err)
}

// TestRenameFileOnExistingFileWithNoReplace tests the renaming of a file on an existing file with no replace.
func (suite *DirectorySuite) TestRenameFileOnExistingFileWithNoReplace() {
	_, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateFile(context.Background(), "newFile", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	err = suite.Directory.RenameFile(context.Background(), "file", suite.Directory, "newFile", true)
	suite.Require().ErrorIs(err, storage.ErrFileAlreadyExists)
}

// TestRenameFileOnExistingFileWithoutNoReplace tests the renaming of a file on
// an existing file without no replace.
func (suite *DirectorySuite) TestRenameFileOnExistingFileWithoutNoReplace() {
	_, err := suite.Directory.CreateFile(context.Background(), "file", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateFile(context.Background(), "newFile", info.File{
		ChunkSize: 8192,
	})
	suite.Require().NoError(err)

	err = suite.Directory.RenameFile(context.Background(), "file", suite.Directory, "newFile", false)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(context.Background(), "file")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	f, err := suite.Directory.GetFile(context.Background(), "newFile")
	suite.Require().NoError(err)

	info, err := f.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(4096, info.ChunkSize)
}

// TestRenameDirectoryOnSameDirectory tests the renaming of a directory on the same directory.
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

// TestRenameDirectoryOnDifferentDirectory tests the renaming of a directory on
// a different directory.
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

// TestRenameDirectoryOnExistingDirectoryWithNoReplace tests the renaming of a
// directory on an existing directory with no replace.
func (suite *DirectorySuite) TestRenameDirectoryOnExistingDirectoryWithNoReplace() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "directory")
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "newDirectory")
	suite.Require().NoError(err)

	err = suite.Directory.RenameDirectory(context.Background(), "directory", suite.Directory, "newDirectory", true)
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}

// TestRenameDirectoryOnExistingDirectoryWithoutNoReplace tests the renaming of
// a directory on an existing directory without no replace.
func (suite *DirectorySuite) TestRenameDirectoryOnExistingDirectoryWithoutNoReplace() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "directory")
	suite.Require().NoError(err)

	_, err = suite.Directory.CreateDirectory(context.Background(), "newDirectory")
	suite.Require().NoError(err)

	err = suite.Directory.RenameDirectory(context.Background(), "directory", suite.Directory, "newDirectory", false)
	suite.Require().NoError(err)

	_, err = suite.Directory.GetDirectory(context.Background(), "directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
}

// TestGetFileWhenIsDirectory tests the GetFile method on a directory.
func (suite *DirectorySuite) TestGetFileWhenIsDirectory() {
	_, err := suite.Directory.CreateDirectory(context.Background(), "dir")
	suite.Require().NoError(err)

	_, err = suite.Directory.GetFile(context.Background(), "dir")
	suite.Require().ErrorIs(err, storage.ErrIsDirectory)
}
