package layer

import (
	"context"

	"github.com/lerenn/chonkfs/pkg/info"
	"github.com/lerenn/chonkfs/pkg/storage"
	"github.com/lerenn/chonkfs/pkg/storage/test"
)

// DirectorySuite is a test suite for directories.
type DirectorySuite struct {
	Upperlayer storage.Directory
	Underlayer storage.Directory
	test.DirectorySuite
}

// TestGetInfoWhenDirectoryExistsOnlyOnUnderlayer tests the retrieval of info
// when a directory only exists on the underlayer.
func (suite *DirectorySuite) TestGetInfoWhenDirectoryExistsOnlyOnUnderlayer() {
	// Create a directory
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Get directory
	dir, err := suite.Directory.GetDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Get info
	dirInfo, err := dir.GetInfo(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(info.Directory{}, dirInfo)
}

// TestRemoveDirectoryOnBackendAndUnderlayer tests the removal of a directory
// when it exists on both the backend and the underlayer.
func (suite *DirectorySuite) TestRemoveDirectoryOnBackendAndUnderlayer() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Remove the directory
	err = suite.Directory.RemoveDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Check it does not exist on directory upperlayer
	_, err = suite.Upperlayer.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)

	// Check it does not exist on underlayer
	_, err = suite.Underlayer.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
}

// TestRemoveDirectoryWhenOnlyOnUnderlayer tests the removal of a directory when it only exists on the underlayer.
func (suite *DirectorySuite) TestRemoveDirectoryWhenOnlyOnUnderlayer() {
	// Create a directory on underlayer
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Remove the directory
	err = suite.Directory.RemoveDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Check it does not exist on directory upperlayer
	_, err = suite.Upperlayer.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)

	// Check it does not exist on underlayer
	_, err = suite.Underlayer.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)
}

// TestListFilesWithOneInUnderlayer tests the listing of files when there is one in the underlayer.
func (suite *DirectorySuite) TestListFilesWithOneInUnderlayer() {
	// Create a directory
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Create a file in underlayer
	_, err = suite.Underlayer.CreateFile(context.Background(), "FileA", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Create 2 files in directory
	_, err = suite.Directory.CreateFile(context.Background(), "FileB", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateFile(context.Background(), "FileC", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// List files
	files, err := suite.Directory.ListFiles(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(files, 3)
}

// TestGetFileWhenOnlyOnUnderlayer tests the retrieval of a file when it only exists on the underlayer.
func (suite *DirectorySuite) TestGetFileWhenOnlyOnUnderlayer() {
	// Create a file on underlayer
	_, err := suite.Underlayer.CreateFile(context.Background(), "File", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Get the file
	file, err := suite.Directory.GetFile(context.Background(), "File")
	suite.Require().NoError(err)
	suite.Require().NotNil(file)
}

// TestGetDirectory tests the retrieval of a directory.
func (suite *DirectorySuite) TestGetDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Get the directory
	directory, err := suite.Directory.GetDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)
	suite.Require().NotNil(directory)
}

// TestGetDirectoryWhenOnlyOnUnderlayer tests the retrieval of a directory when it only exists on the underlayer.
func (suite *DirectorySuite) TestGetDirectoryWhenOnlyOnUnderlayer() {
	// Create a directory on underlayer
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Get the directory
	directory, err := suite.Directory.GetDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)
	suite.Require().NotNil(directory)
}

// TestCreateDirectory tests the creation of a directory.
func (suite *DirectorySuite) TestCreateDirectory() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Check it exists on directory upperlayer
	_, err = suite.Upperlayer.GetDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Check it exists on underlayer
	_, err = suite.Underlayer.GetDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)
}

// TestCreateDirectoryWhenDirectoryAlreadyExists tests the creation of a directory
// when a directory with the same name already exists.
func (suite *DirectorySuite) TestCreateDirectoryWhenDirectoryAlreadyExists() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Create the same directory again
	_, err = suite.Directory.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().ErrorIs(err, storage.ErrDirectoryAlreadyExists)
}

// TestCreateDirectoryWhenFileAlreadyExists tests the creation of a directory
// when a file with the same name already exists.
func (suite *DirectorySuite) TestCreateDirectoryWhenFileAlreadyExists() {
	// Create a file
	_, err := suite.Directory.CreateFile(context.Background(), "test", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Create a directory with the same name
	_, err = suite.Directory.CreateDirectory(context.Background(), "test")
	suite.Require().ErrorIs(err, storage.ErrFileAlreadyExists)
}

// TestListDirectoriesWithOneInUnderlayer tests the listing of directories when there is one in the underlayer.
func (suite *DirectorySuite) TestListDirectoriesWithOneInUnderlayer() {
	// Create a file
	_, err := suite.Underlayer.CreateFile(context.Background(), "File", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Create a file in underlayer
	_, err = suite.Underlayer.CreateDirectory(context.Background(), "DirectoryA")
	suite.Require().NoError(err)

	// Create 2 files in directory
	_, err = suite.Directory.CreateDirectory(context.Background(), "DirectoryB")
	suite.Require().NoError(err)
	_, err = suite.Directory.CreateDirectory(context.Background(), "DirectoryC")
	suite.Require().NoError(err)

	// List directories
	dirs, err := suite.Directory.ListDirectories(context.Background())
	suite.Require().NoError(err)
	suite.Require().Len(dirs, 3)
}

// TestRemoveFileOnBackendAndUnderlayer tests the removal of a file when it
// exists on both the backend and the underlayer.
func (suite *DirectorySuite) TestRemoveFileOnBackendAndUnderlayer() {
	// Create a directory
	_, err := suite.Directory.CreateFile(context.Background(), "File", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Remove the directory
	err = suite.Directory.RemoveFile(context.Background(), "File")
	suite.Require().NoError(err)

	// Check it does not exist on directory upperlayer
	_, err = suite.Upperlayer.GetFile(context.Background(), "File")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	// Check it does not exist on underlayer
	_, err = suite.Underlayer.GetFile(context.Background(), "File")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)
}

// TestRemoveFileWhenOnlyOnUnderlayer tests the removal of a file when it only exists on the underlayer.
func (suite *DirectorySuite) TestRemoveFileWhenOnlyOnUnderlayer() {
	// Create a directory on underlayer
	_, err := suite.Underlayer.CreateFile(context.Background(), "File", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Remove the directory
	err = suite.Directory.RemoveFile(context.Background(), "File")
	suite.Require().NoError(err)

	// Check it does not exist on directory upperlayer
	_, err = suite.Upperlayer.GetFile(context.Background(), "File")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	// Check it does not exist on underlayer
	_, err = suite.Underlayer.GetFile(context.Background(), "File")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)
}

// TestRenameFileOnBackendAndUnderlayer tests the renaming of a file when it
// exists on both the backend and the underlayer.
func (suite *DirectorySuite) TestRenameFileOnBackendAndUnderlayer() {
	// Create a directory
	_, err := suite.Directory.CreateFile(context.Background(), "File", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Rename the directory
	err = suite.Directory.RenameFile(context.Background(), "File", suite.Directory, "File2", false)
	suite.Require().NoError(err)

	// Check it does not exist on directory upperlayer
	_, err = suite.Upperlayer.GetFile(context.Background(), "File")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	// Check it does not exist on underlayer
	_, err = suite.Underlayer.GetFile(context.Background(), "File")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	// Check it exists on directory upperlayer
	_, err = suite.Upperlayer.GetFile(context.Background(), "File2")
	suite.Require().NoError(err)

	// Check it exists on underlayer
	_, err = suite.Underlayer.GetFile(context.Background(), "File2")
	suite.Require().NoError(err)
}

// TestRenameFileWhenOnlyOnUnderlayer tests the renaming of a file when it only exists on the underlayer.
func (suite *DirectorySuite) TestRenameFileWhenOnlyOnUnderlayer() {
	// Create a directory on underlayer
	_, err := suite.Underlayer.CreateFile(context.Background(), "File", info.File{
		ChunkSize: 4096,
	})
	suite.Require().NoError(err)

	// Rename the directory
	err = suite.Directory.RenameFile(context.Background(), "File", suite.Directory, "File2", false)
	suite.Require().NoError(err)

	// Check it does not exist on directory
	_, err = suite.Directory.GetFile(context.Background(), "File")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	// Check it does not exist on underlayer
	_, err = suite.Underlayer.GetFile(context.Background(), "File")
	suite.Require().ErrorIs(err, storage.ErrFileNotFound)

	// Check it exists on directory
	_, err = suite.Directory.GetFile(context.Background(), "File2")
	suite.Require().NoError(err)

	// Check it exists on underlayer
	_, err = suite.Underlayer.GetFile(context.Background(), "File2")
	suite.Require().NoError(err)
}

// TestRenameOnBackendAndUnderlayer tests the renaming of a directory when it
// exists on both the backend and the underlayer.
func (suite *DirectorySuite) TestRenameOnBackendAndUnderlayer() {
	// Create a directory
	_, err := suite.Directory.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Rename the directory
	err = suite.Directory.RenameDirectory(context.Background(), "Directory", suite.Directory, "Directory2", false)
	suite.Require().NoError(err)

	// Check it does not exist on directory upperlayer
	_, err = suite.Upperlayer.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)

	// Check it does not exist on underlayer
	_, err = suite.Underlayer.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)

	// Check it exists on directory upperlayer
	_, err = suite.Upperlayer.GetDirectory(context.Background(), "Directory2")
	suite.Require().NoError(err)

	// Check it exists on underlayer
	_, err = suite.Underlayer.GetDirectory(context.Background(), "Directory2")
	suite.Require().NoError(err)
}

// TestRenameDirectoryWhenOnlyOnUnderlayer tests the renaming of a directory when it only exists on the underlayer.
func (suite *DirectorySuite) TestRenameDirectoryWhenOnlyOnUnderlayer() {
	// Create a directory on underlayer
	_, err := suite.Underlayer.CreateDirectory(context.Background(), "Directory")
	suite.Require().NoError(err)

	// Rename the directory
	err = suite.Directory.RenameDirectory(context.Background(), "Directory", suite.Directory, "Directory2", false)
	suite.Require().NoError(err)

	// Check it does not exist on directory
	_, err = suite.Directory.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)

	// Check it does not exist on underlayer
	_, err = suite.Underlayer.GetDirectory(context.Background(), "Directory")
	suite.Require().ErrorIs(err, storage.ErrDirectoryNotFound)

	// Check it exists on directory
	_, err = suite.Directory.GetDirectory(context.Background(), "Directory2")
	suite.Require().NoError(err)

	// Check it exists on underlayer
	_, err = suite.Underlayer.GetDirectory(context.Background(), "Directory2")
	suite.Require().NoError(err)
}
